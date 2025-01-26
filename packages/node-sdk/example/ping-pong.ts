import { ToolProvider } from '../src/platform-provider'

const tool = new ToolProvider({
    project: 'sandbox',
    endpoint: 'ws://localhost:8080/ws/tool/provider',
    credentials: {
        username: 'test@test.com',
        password: 'password',
    },
    tools: [
        {
            toolId: 'sandbox',
            toolFunction: async ({ io }) => {
                const input = await io.prompt(
                    'Input ping to receive pong',
                    'string'
                )

                if (input.toLowerCase() === 'ping') {
                    return 'pong'
                }

                throw '...'
            },
        },
    ],
})

tool.listen()
