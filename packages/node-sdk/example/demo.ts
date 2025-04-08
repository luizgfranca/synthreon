import { ToolProvider } from '../src/tool-provider'

const provider = new ToolProvider({
    project: 'sandbox',
    endpoint: 'ws://localhost:25256/ws/tool/provider',
    credentials: {
        username: 'test@test.com',
        password: 'password',
    },
})

provider.tool('tool-y', async ({ io }) => {
    const input = await io.prompt({
        title: 'Try to input a ping command',
        type: 'string'
    })

    if (input.toLowerCase() === 'ping') {
        return 'pong';
    }

    throw 'You failed.'
}, {
    name: 'Test Y',
    description: 'This is the Y tool test'
})

provider.listen()
