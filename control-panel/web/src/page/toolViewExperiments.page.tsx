import { DisplayDefinition, DisplayRenderer, DsiplayRendererProps } from "@/component/displayRenderer"
import { useCallback, useMemo, useState } from "react"

// const event = {
//     "class": "operation",
//     "type": "display",
//     "project": "proj-x",
//     "tool": "tool-y",
//     "display": {
//         "type": "result",
//         "result": {
//             "success": true,
//             "message": "Hello user input",
//         }
//     }
// }

type ToolEvent = {
    "class": string;
    "type": string;
    "project": string;
    "tool": string;
    "display": DisplayDefinition;
}

export function ToolViewExperimentsPage() {
    const [event, setEvent] = useState<ToolEvent | null>(null)

    console.log('e', event)
    
    useMemo(() => {
        const ws = new WebSocket(`${import.meta.env.PL_BACKEND_URL}/api/tool/client/ws`)
        ws.addEventListener('open', () => {
            console.log('socket open')

            ws.send(JSON.stringify({
                "class": "interaction",
                "type": "open",
                "project": "proj-x",
                "tool": "tool-b",
            }))
        })

        ws.addEventListener('message', (e) => {
            console.log(`recv: ${e.data}`)
            setEvent(JSON.parse(e.data))
        })

        return ws
    }, [])
    
    if(!event) {
        return (
            <div className="bg-zinc-900 text-zinc-100 h-screen">
                    <div className="container mx-auto px-4 py-8">
                    <h1 className="text-3xl font-bold mb-6">Tool Sandbox</h1>
                </div>
            </div>
        )       
    }
    return (
        <div className="bg-zinc-900 text-zinc-100 h-screen">
            <div className="container mx-auto px-4 py-8">
                <h1 className="text-3xl font-bold mb-6">Tool Sandbox</h1>
                <DisplayRenderer definition={event.display}/>
            </div>
        </div>
    )
}