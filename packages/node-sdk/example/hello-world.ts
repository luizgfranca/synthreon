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
            function: async () => {
                return 'Hello World!'
            },
        },
    ],
})

tool.listen()
