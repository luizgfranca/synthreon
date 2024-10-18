import WebSocket from 'ws'


const ws = new WebSocket('ws://localhost:8080/api/tool/provider/ws')

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
    console.log('event', e.toString())
})

ws.on('error', (error) => {
    console.log('error', error)
})