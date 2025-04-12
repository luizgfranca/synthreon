import useWebSocket, { ReadyState } from 'react-use-websocket'
import { EmptyState } from "@/component/emptyState";
import Compositor from "@/component/compositor";
import BackendService from "@/service/backend.service";
import { useCallback, useEffect, useState } from "react";
import { ToolEventDto } from '@synthreon/core';

const ToolEventEncoder = {
    encodeV0: (event: ToolEventDto): string => {
        const prefix = 'v0.0';
        const data = JSON.stringify(event)
        return `${prefix}|${data}`
    },
    decodeV0: (input: string): ToolEventDto => {
        const [prefix, rawData] = input.split('|')
        if (prefix !== 'v0.0') {
            throw new Error('unsupported protovol version')
        }

        return JSON.parse(rawData)
    }
}


type ToolViewProps = {
    project: string,
    tool: string
}

enum ToolState {
    STARTUP = 'STARTUP',
    OPENNING = 'OPENNING',
    OPEN = 'OPEN',
}

type ExecutionContext = {
    contextId?: string;
    terminalId?: string;
    state: ToolState
}

function empty(message: string) {
    return (
        <div className="flex pt-10 items-center">
            <EmptyState>{message}</EmptyState>
        </div>
    )
}

// dont need to verify if project and tool are different because of the key property
// that will reset the component if they are
const ToolView = (props: ToolViewProps) => {
    const [context, setContext] = useState<ExecutionContext>({ state: ToolState.STARTUP })
    const [eventHistory, setEventHistory] = useState<ToolEventDto[]>([]);
    console.debug('rendering tool with context', context, props)

    const reset = useCallback(() => {
        console.debug('reset')
        setEventHistory([]);
        setContext({ state: ToolState.STARTUP, })
    }, [setContext])

    const accessToken = BackendService.getAccessToken() ?? ''
    const BASE_URL = `http://${window.location.hostname}:25256`
    const { sendMessage, lastMessage, readyState } = useWebSocket(`${BASE_URL}/api/tool/client/ws/${accessToken}`);

    const sendEvent = useCallback((toSend: ToolEventDto) => {
        const event: ToolEventDto = {
            ...toSend,
            project: props.project,
            tool: props.tool
        }

        // TODO: could validate event here in the future for sanity checking
        console.debug('sending event:', event);
        const message = ToolEventEncoder.encodeV0(event);
        sendMessage(message);
    }, [])

    useEffect(() => {
        if (lastMessage !== null) {
            let lastEvent: ToolEventDto;
            try {
                console.debug('message received', lastMessage);
                // TODO: validate received event
                lastEvent = ToolEventEncoder.decodeV0(lastMessage.data)
                setEventHistory((prev) => prev.concat(lastEvent));
            } catch (e) {
                // FIXME: error screen for this
                console.error('unsupported event received', e);
            }
        }
    }, [lastMessage]);

    let lastEvent = eventHistory[eventHistory.length - 1];
    switch (context.state) {
        case ToolState.STARTUP:
            if (readyState == ReadyState.OPEN) {
                setContext({ state: ToolState.OPENNING })
                sendEvent({ type: 'interaction/open' })
                // FIXME: loading screen for this
                return empty('Starting tool...')
            }
            // FIXME: loading screen for this
            return empty('Connecting...')
        case ToolState.OPENNING:
            if (!lastEvent) {
                return empty('Waiting response...');
            }

            try {
                // TODO: validate received event
                setContext({
                    contextId: lastEvent.context_id,
                    terminalId: lastEvent.terminal_id,
                    state: ToolState.OPEN
                })

                // FIXME: loading screen for this
                return empty('Openning...');
            } catch (e) {
                // FIXME: error screen for this
                return empty('Unsupported event received...');
            }
        case ToolState.OPEN:
            if (!lastEvent) {
                console.error('state is OPEN but last event is not here yet', props, context);
                return empty('internal error, invalid state ...');
            }

            if (context.contextId && context.terminalId) {
                console.debug('loaded', { ...props, ...context });
                return <Compositor
                    key={`${context.contextId}:${context.terminalId}`}
                    contextId={context.contextId}
                    terminalId={context.terminalId}
                    lastEvent={lastEvent}
                    reset={reset}
                    sendEvent={sendEvent}
                />
            }
            return empty('should be unreacheable')
    }
}


export default ToolView;
