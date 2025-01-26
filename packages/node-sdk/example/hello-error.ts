import { ToolProvider } from '../src/platform-provider'

const tool = new ToolProvider({
    project: 'sandbox',
    endpoint: 'ws://localhost:8080/api/tool/provider/ws',
    credentials: {
        username: 'test@test.com',
        password: 'password',
    },
    tools: [
        {
            toolId: 'sandbox',
            toolFunction: async () => {
                throw 'Hello error'
            },
        },
    ],
})

tool.listen()
