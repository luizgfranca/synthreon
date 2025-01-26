import WebSocket from 'ws'


const ws = new WebSocket('ws://localhost:8080/api/tool/provider/ws')


type DisplayElement = {
    type: string;
    label: string;
    text: string;
    description: string;
    name: string;
};

type DisplayResult = {
    success: boolean;
    message: string;
};

type DisplayDefinitionType = 'result' | 'view' | 'prompt' | string

type DisplayPrompt = {
    title: string;
    type: string;
}

export type DisplayDefinition = {
    type: DisplayDefinitionType;
    elements?: DisplayElement[];
    result?: DisplayResult;
    prompt?: DisplayPrompt;
};

export type DsiplayRendererProps = {
    definition: DisplayDefinition
}

export type InputField = {
    name: string;
    value: string;
}

export type EventInput = {
    fields: InputField[]
}

type ToolEvent = {
    class: string;
    type: string;
    project: string;
    tool: string;
    client?: string;
    display?: DisplayDefinition;
    input?: EventInput
}

ws.on('open', () => {
    console.log('opened')
    ws.send(JSON.stringify({
        "class": "announcement",
        "type": "provider",
        "project": "sandbox",
        "tool": "sandbox",
        "provider": 1
    }))
})

ws.on('message', (e) => {
    const event = JSON.parse(e.toString()) as ToolEvent
    console.log('event', event)

    if(event.class == 'interaction' && event.type == 'open' && event.project == 'proj-x') {
        if(event.tool == 'tool-y') {
            const response: ToolEvent = {
                class: 'operation',
                type: 'display',
                project: event.project,
                tool: event.tool,
                client: event.client,
                display: {
                    type: 'result',
                    result: {
                        success: true,
                        message: 'success running tool A'
                    }
                }
            }

            console.log('response event', response)
            ws.send(JSON.stringify(response))
        } else if(event.tool == 'tool-p') {
            const response: ToolEvent = {
                class: 'operation',
                type: 'display',
                project: event.project,
                tool: event.tool,
                client: event.client,
                display: {
                    type: 'prompt',
                    prompt: {
                        title: 'Add some text in the field',
                        type: 'string'
                    }
                }
            }

            console.log('response event', response)
            ws.send(JSON.stringify(response))
        } else {
            let incrementString = ''

            setInterval(() => {
                console.log('settimeout loop')
                incrementString += 'x';

                const response: ToolEvent = {
                    class: 'operation',
                    type: 'display',
                    project: event.project,
                    tool: event.tool,
                    client: event.client,
                    display: {
                        type: 'result',
                        result: {
                            success: true,
                            message: 'success running alternative tool' + incrementString
                        }
                    }
                }
    
                console.log('response event', response)
                ws.send(JSON.stringify(response))    
            }, 1000)
        }   
    } else if(event.class == 'interaction' && event.type == 'input' && event.project == 'proj-x') {
        console.log('input interaction', event)
        
        const response: ToolEvent = {
            class: 'operation',
            type: 'display',
            project: event.project,
            tool: event.tool,
            client: event.client,
            display: {
                type: 'result',
                result: {
                    success: true,
                    message: event.input?.fields[0].value ?? ''
                }
            }
        }

        console.log('response event', response)
        ws.send(JSON.stringify(response))      
    }
})

ws.on('error', (error) => {
    console.log('error', error)
})