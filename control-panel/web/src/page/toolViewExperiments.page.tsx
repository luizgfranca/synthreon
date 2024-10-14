import { DisplayRenderer } from "@/component/displayRenderer"

const event = {
    "class": "operation",
    "type": "display",
    "project": "proj-x",
    "tool": "tool-y",
    "display": {
        "type": "result",
        "result": {
            "success": true,
            "message": "Hello user input",
        }
    }
}



export function ToolViewExperimentsPage() {
    return (
        <div className="bg-zinc-900 text-zinc-100 h-screen">
            <div className="container mx-auto px-4 py-8">
                <h1 className="text-3xl font-bold mb-6">Tool Sandbox</h1>
                <DisplayRenderer definition={event.display}/>
            </div>
        </div>
    )
}