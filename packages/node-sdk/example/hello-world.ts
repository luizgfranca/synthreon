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
            toolFunction: async () => {
                return 'Hello World!'
            },
        },
    ],
})

tool.listen()
