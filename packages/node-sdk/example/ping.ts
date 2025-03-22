import { ToolProvider } from '../src/tool-provider'

const provider = new ToolProvider({
    project: 'sandbox',
    endpoint: 'ws://localhost:25256/ws/tool/provider',
    credentials: {
        username: 'test@test.com',
        password: 'password',
    },
})

provider.tool('sandbox', async ({ io }) => {
    const input = await io.prompt({
        title: 'Try to input a ping command',
        type: 'string'
    })

    if (input.toLowerCase() === 'ping') {
        return 'pong';
    }

    throw 'You failed.'
})

provider.listen()
