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

export type DisplayDefinition = {
    type: DisplayDefinitionType;
    elements?: DisplayElement[];
    result?: DisplayResult;
};

export type DsiplayRendererProps = {
    definition: DisplayDefinition
}

type ToolEvent = {
    class: string;
    type: string;
    project: string;
    tool: string;
    client?: string;
    display: DisplayDefinition;
}

ws.on('open', () => {
    console.log('opened')
    ws.send(JSON.stringify({
        "class": "announcement",
        "type": "provider",
        "project": "proj-x",
        "tool": "tool-y",
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
                type: 'result',
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
        } else {
            let incrementString = ''

            setInterval(() => {
                console.log('settimeout loop')
                incrementString += 'x';

                const response: ToolEvent = {
                    class: 'operation',
                    type: 'result',
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
    }
})

ws.on('error', (error) => {
    console.log('error', error)
})