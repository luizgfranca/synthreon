import { ToolProvider } from '../src/tool-provider'

const tool = new ToolProvider({
    project: 'sandbox',
    endpoint: 'ws://localhost:25256/ws/tool/provider',
    credentials: {
        username: 'test@test.com',
        password: 'password',
    },
    tools: [
        {
            id: 'sandbox',
            function: async ({ io }) => {
                const input = await io.prompt(
                    {
                        title: 'Simple ping test',
                        type: 'string'
                    }
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
