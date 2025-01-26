import { DisplayRenderer, Field } from "@/component/displayRenderer";
import { EmptyState } from "@/component/emptyState";
import BackendService from "@/service/backend.service";
import { useMemo, useState } from "react";

import { ToolEventDto } from 'platformlab-core';


const ToolEventEncoder = {
    encodeV0: (event: ToolEventDto): string => {
        const prefix = 'v0.0';
        const data = JSON.stringify(event)
        return `${prefix}|${data}`
    },
    decodeV0: (input: string): {result?: ToolEventDto, error?: Error} => {
        const [prefix, rawData] = input.split('|')
        if(prefix !== 'v0.0') {
            return {
                error: new Error('unsupported protovol version')
            }
        }

        return {
            result: JSON.parse(rawData)
        }
    }
}

type ToolViewProps = {
    project?: string,
    tool?: string,
}

type ExecutionContext = {
    contextId?: string;
    terminalId?: string;
}

function sendEvent(ws: WebSocket, event: ToolEventDto) {
    console.log('sending event: ', event);
    ws.send(ToolEventEncoder.encodeV0(event))
}

export function ToolView(props: ToolViewProps) {
    const [event, setEvent] = useState<ToolEventDto | null>(null)
    const [executionContext, setExecutionContext] = useState<ExecutionContext | undefined>()
    const [resetToggle, setResetToggle] = useState<boolean>()
    console.log('e', event)

    const accessToken = BackendService.getAccessToken() ?? ''
    
    const BASE_URL = `http://${window.location.hostname}:8080`
    const ws = useMemo(() => {
        const ws = new WebSocket(
            `${BASE_URL}/api/tool/client/ws`,
            [ accessToken ]
        )
        ws.addEventListener('open', () => {
            console.log('socket open', {
                project: props.project,
                tool: props.tool
            })

            if (props.project && props.tool) {
                sendEvent(ws, {
                    type: 'interaction/open',
                    project: props.project,
                    tool: props.tool
                })
            } else {
                throw new Error('project ant tool should be defined to send an interaction')
            }
        })

        ws.addEventListener('message', (e) => {
            console.log(`recv: ${e.data}`)
            const {result, error} = ToolEventEncoder.decodeV0(e.data)
            if (!error && result) {
                if(!executionContext) {
                    console.debug('setting execution context', {
                        contextId: result.context_id,
                        terminalId: result.terminal_id
                    })

                    setExecutionContext({
                        contextId: result.context_id,
                        terminalId: result.terminal_id
                    })
                }
                setEvent(result)
            } else {
                console.error('error decoding event: ', error)
            }
        })

        setResetToggle(false)

        return ws
    }, [props, resetToggle])
 
    const sendInputInteraction = (fields: Field[]) => {
        if (props.project && props.tool) {
            sendEvent(ws, {
                type: 'interaction/input',
                project: props.project,
                tool: props.tool,
                context_id: executionContext?.contextId,
                terminal_id: executionContext?.terminalId,
                input: {
                    fields
                }
            })
        } else {
            throw new Error('project ant tool should be defined to send an interaction')
        }
    }

    const reset = () => {
        console.log('RESET')
        setResetToggle(true)
        setExecutionContext(undefined)
    }

    if (!props.project || !props.tool) {
        return (
            <div className="flex pt-10 items-center">
                <EmptyState>Select a tool to use.</EmptyState>
            </div>
        )
    }

    if(!event || (!event.display && !event.result)) {
        console.debug('no event', event)
        return (
            <div className="flex pt-10">

                <EmptyState>Waiting for provider...</EmptyState>
            </div>
        )       
    }
    return (
        <div className="text-zinc-100 h-screen pt-10">
            <div className="container mx-auto px-4">
                <DisplayRenderer 
                    event={event}
                    onSumission={(fields) => sendInputInteraction(fields)}
                    resetCallback={() => reset()}
                />
            </div>
        </div>
    )
}