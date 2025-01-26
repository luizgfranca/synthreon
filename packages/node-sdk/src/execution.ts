import { EventEmitter } from 'node:events'
import { ToolComponents, ToolHandlerDefinition } from './handler'
import { ComponentFactory } from './tool-components'
import { EventTypeValue, ToolEventDto } from 'platformlab-core'
import { v4 as uuid } from 'uuid'
import { InputDefinition } from 'platformlab-core/tool-event/input/input.dto'
import { ToolEventResult } from 'platformlab-core/tool-event/result/result.dto'
import { DisplayDefinition } from 'platformlab-core'

type ForwardEventCallback = (event: ToolEventDto) => void
type InputSubscriber = (input: InputDefinition) => void

export class Execution {
    #id: string
    #contextId?: string

    // FIXME: subscribers should have a timeout
    #inputSubscribers: InputSubscriber[]

    #bus: EventEmitter
    #components: ToolComponents
    #definition: ToolHandlerDefinition
    #promise?: Promise<string>
    #forwardEvent: ForwardEventCallback

    constructor(
        bus: EventEmitter, 
        definition: ToolHandlerDefinition,
        forwardEventCallback: ForwardEventCallback
    ) {
        this.#id = uuid();
        this.#bus = bus;
        this.#definition = definition;
        this.#forwardEvent = forwardEventCallback;
        
        this.#components = ComponentFactory.instantiateComponents(this);
        this.#inputSubscribers = []
    }

    start(startingEvent: ToolEventDto) {
        this.#promise = new Promise((resolve, reject) => {
            try {
                this.#definition
                    .toolFunction(this.#components)
                    .then((resultMessage) => resolve(resultMessage))
                    .catch((e) => reject(e))
            } catch (e) {
                reject(e)
            }
        });

        this.#contextId = startingEvent.context_id;

        this.#promise
            .then((resultMessage) => {
                console.debug('executing success result trap')
                this.#bus.removeAllListeners(this.#id)
                this.sendResult({status: 'success', message: resultMessage})
            })
            .catch((errorMessage) => {
                console.debug('executing error result trap')
                this.#bus.removeAllListeners(this.#id)
                this.sendResult({status: 'success', message: errorMessage})
            })
        
        this.#bus.on(
            this.#id, 
            (e: ToolEventDto) => this.#onEvent(e)
        )
    }

    onNextInput(subscriber: InputSubscriber) {
        this.#inputSubscribers.push(subscriber)
    }

    sendDisplay(definition: DisplayDefinition) {
        this.#sendEvent({
            type: EventTypeValue.CommandDisplay,
            display: definition
        })
    }

    sendResult(result: ToolEventResult) {
        this.#sendEvent({
            type: EventTypeValue.CommandFinish,
            result
        })
    }

    #onEvent(event: ToolEventDto) {
        switch(event.type) {
            case EventTypeValue.InteractionInput:
                this.#onInput(event);
                return;
            default:
                console.warn('DROPPING: invalid event type to be sent to execution', event)
        }
    }

    #onInput(event: ToolEventDto) {
        console.debug('received input event', event)

        if(!event.input) {
            console.warn('received input event with no input defined');
            return;
        }

        if(!this.#inputSubscribers.length) {
            console.warn('no subscribers to receive interaction');
        }

        this.#inputSubscribers.forEach(subscriber => {
            if(!event.input) {
                console.error('internal error: input state not available on subscriber execution')
                process.exit()
            }

            subscriber(event.input)
        });

        this.#inputSubscribers = []
    }

    #sendEvent(event: ToolEventDto) {
        event.tool = this.#definition.toolId
        event.execution_id = this.#id
        event.context_id = this.#contextId

        this.#forwardEvent(event)
    }
}
