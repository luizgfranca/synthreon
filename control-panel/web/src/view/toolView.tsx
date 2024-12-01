import { DisplayRenderer, Field } from "@/component/displayRenderer";
import { EmptyState } from "@/component/emptyState";
import { ToolEvent } from "@/model/toolEvent";
import BackendService from "@/service/backend.service";
import { useMemo, useState } from "react";

type ToolViewProps = {
    project?: string,
    tool?: string,
}

function sendEvent(ws: WebSocket, event: ToolEvent) {
    console.log('sending event: ', event);
    ws.send(JSON.stringify(event))
}

export function ToolView(props: ToolViewProps) {
    const [event, setEvent] = useState<ToolEvent | null>(null)
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
                    class: 'interaction',
                    type: 'open',
                    project: props.project,
                    tool: props.tool
                })
            } else {
                throw new Error('project ant tool should be defined to send an interaction')
            }
        })

        ws.addEventListener('message', (e) => {
            console.log(`recv: ${e.data}`)
            setEvent(JSON.parse(e.data))
        })

        return ws
    }, [props])
 
    const sendInputInteraction = (fields: Field[]) => {
        if (props.project && props.tool) {
            sendEvent(ws, {
                class: 'interaction',
                type: 'input',
                project: props.project,
                tool: props.tool,
                input: {
                    fields
                }
            })
        } else {
            throw new Error('project ant tool should be defined to send an interaction')
        }
    }

    if (!props.project || !props.tool) {
        return (
            <div className="flex pt-10 items-center">
                <EmptyState>Select a tool to use.</EmptyState>
            </div>
        )
    }

    if(!event || !event.display) {
        return (
            <div className="flex pt-10">

                <EmptyState>Waiting for provider...</EmptyState>
            </div>
        )       
    }
    return (
        <div className="text-zinc-100 h-screen pt-10">
            <div className="container mx-auto px-4">
                <DisplayRenderer definition={event.display} onSumission={(fields) => sendInputInteraction(fields)}/>
            </div>
        </div>
    )
}