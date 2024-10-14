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
    
    useMemo(() => {
        const s = new WebSocket(`${import.meta.env.PL_BACKEND_URL}/api/tool/ws`)
        s.addEventListener('open', () => {
            console.log('socket open')
            setInterval(() => {
                s.send(JSON.stringify({
                    "class": "interaction",
                    "type": "open",
                    "project": "proj-x",
                    "tool": "tool-y",
                }))
                console.log('sent')
            }, 1000)
        })

        s.addEventListener('message', (e) => {
            console.log(`recv: ${e.data}`)
            setEvent(JSON.parse(e.data))
        })

        return s
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

    console.log('e', event)

    return (
        <div className="bg-zinc-900 text-zinc-100 h-screen">
            <div className="container mx-auto px-4 py-8">
                <h1 className="text-3xl font-bold mb-6">Tool Sandbox</h1>
                <DisplayRenderer definition={event.display}/>
            </div>
        </div>
    )
}