import { PromptTypeOption, PromptType, ToolEventDto, EventTypeValue, ToolEventEncoder } from 'platformlab-core'
import { EventEmitter } from 'node:events'
import WebSocket, { RawData } from 'ws'
import { Handler, HandlerExtraOptions, ToolFunction, ToolHandlerDefinition } from './handler'

type UserCredentials = {
    username: string
    password: string
}

type ToolOptions = HandlerExtraOptions;

type PlatformToolConnectionOptions = {
    endpoint: string,
    project: string,
    credentials: UserCredentials,
    tools?: ToolHandlerDefinition[]
}

type ConnectionStatus = 'connected' 
    | 'connecting' 
    | 'waiting_ack' 
    | 'disconnected'

export class ToolProvider {
    #endpoint: string
    #websocket?: WebSocket
    #status: ConnectionStatus
    #credentials: UserCredentials
    #project: string

    #providerId?: string
    #handshakeId?: string

    #bus?: EventEmitter;

    #handlerDefinitions: ToolHandlerDefinition[]
    #handlers: Handler[]

    constructor(options: PlatformToolConnectionOptions) {
        this.#endpoint = options.endpoint;
        this.#credentials = options.credentials;
        this.#status = 'disconnected';
        this.#project = options.project;
        this.#handlerDefinitions = options.tools ? options.tools : [];
        this.#handlers = [];
    }

    // FIXME: think of the better way to allow the definition of a name and a description
    //        for the tool here, mainly for when it has to be autocreated 
    tool(id: string, fn: ToolFunction, options?: ToolOptions) {
        this.#handlerDefinitions.push({
            id,
            function: fn,
            extraOptions: options ? options : {}
        });
    }

    listen() {
        console.log('authstr', `${this.#credentials.username}:${this.#credentials.password}`)
        console.log('authorization', `Basic ${Buffer.from(
                        `${this.#credentials.username}:${
                            this.#credentials.password
                        }`
                    ).toString('base64')}`)

        try {
            this.#websocket = new WebSocket(this.#endpoint, {
                headers: {
                    Authorization: `Basic ${Buffer.from(
                        `${this.#credentials.username}:${
                            this.#credentials.password
                        }`
                    ).toString('base64')}`,
                },
            })
            this.#status = 'connecting'
        } catch (e) {
            throw new Error(`unable to connect to websocket: ${e}`)
        }

        this.#websocket.on('open', () => {
            console.debug('successfully connected to platform')
            this.#requestHandshake()
        })

        this.#websocket.on("message", (data: RawData) => {
            const { result, error } = ToolEventEncoder.decodeV0(data.toString());
            if (error) {
                console.error('unable to parse event', data.toString())
            }
            if (!result) {
                console.error('internal error: or a error or a result should be defined')
                return;
            }

            this.#onEventReceived(result);
        });

    }

    #onEventReceived(event: ToolEventDto) {
        switch(this.#status) {
            case 'waiting_ack':
                this.#processHandshakeResponse(event)
                break;
            case 'connected':
                this.#processNormalEvent(event)
                break;
            default:
                throw new Error('unexpected state: message received while status is disconnected or connecting')
        }
    }

    #sendEvent(event: ToolEventDto) {
        console.debug('sending event: ', event)

        if (!this.#websocket) {
            throw new Error('Internal error: expected websocket to be defined.')
        }
        this.#websocket.send(ToolEventEncoder.encodeV0(event))
    }

    #processNormalEvent(event: ToolEventDto) {
        if(!this.#bus) {
            throw new Error('Internal invalid state: receiving event as connected but bus not set up properly')
        }

        if(!event.project || !event.tool) {
            console.warn('DROPPING malformed event', event)
            return
        }

        if (
            event.type == EventTypeValue.AnnouncementACK
            || event.type == EventTypeValue.AnnouncementNACK
        ) {
            this.#bus.emit(`announcement/${event.tool}`, event)
            return
        }

        if(event.type === EventTypeValue.InteractionOpen) {
            console.debug('open event, directing to handler')
            if (!event.handler_id) {
                console.warn('DROPPING: open event without handler_id', event)
                return
            }
            
            this.#bus.emit(event.handler_id, event)
            return
        }

        if(!event.execution_id) {
            console.warn('DROPPING: non announcement event without execution_id', event)
            return
        }

        // FIXME: direct access like this is dangerous. Maybe I should add some validations here?
        this.#bus.emit(event.execution_id, event)
    }

    #onBusSendRequest(e: ToolEventDto) {
        e.handshake_id = this.#handshakeId;
        e.provider_id = this.#providerId;
        e.project = this.#project;
        this.#sendEvent(e);
    }

    #setupBus() {
        this.#bus = new EventEmitter()
        this.#bus.on('send', (e) => this.#onBusSendRequest(e))
    }

    #getHandshakeEvent(): ToolEventDto {
        return {
            type: EventTypeValue.HandshakeRequest,
            project: this.#project,
        }
    }

    #requestHandshake() {
        if (!this.#websocket) {
            throw new Error(
                'Internal error: expected websocket to be defined to start handshake'
            )
        }

        console.debug("requesting provider handshake");
        this.#sendEvent(this.#getHandshakeEvent());
        this.#status = "waiting_ack";

        console.debug("handshke request sent");
    }

    #processHandshakeResponse(e: ToolEventDto) {
        console.debug("handling expected handshake response");
        if (!e || !e.type) {
            throw new Error('invalid event format')
        }
        
        switch(e.type) {
            case EventTypeValue.HandshakeNACK:
                throw new Error(`handshake request not accepted by server; reason: ${e.reason}`)
            case EventTypeValue.HandshakeACK:
                this.#providerId = e.provider_id
                this.#handshakeId = e.handshake_id
                this.#setupBus()
                this.#status = 'connected'
                this.#startAnnouncement()
                return;
            default:
                throw new Error(`invalid event type received for handshake: ${e.type}`)
        }
    }

    #startAnnouncement() {
        console.debug('starting announcement')
        // TODO: evaluate if i should pass the bus on the start function to be able to create handlers on the constructor
        this.#handlers = this.#handlerDefinitions.map(
            definition => {
                if (!this.#bus) {
                    throw new Error('internal error: trying to do announcement without an active bus')
                }

                console.debug('creating handler:', definition.id)
                return new Handler(definition, this.#bus)
            }
        )

        this.#handlers.forEach(handler => {
            handler.start()
        });
    }
}
