import { ToolProvider } from '../src/tool-provider'

const tool = new ToolProvider({
    project: 'sandbox',
    endpoint: 'ws://localhost:8080/ws/tool/provider',
    credentials: {
        username: 'test@test.com',
        password: 'password',
    },
    tools: [
        {
            id: 'sandbox',
            function: async ({ io }) => {
                const selected = await io.selection({
                    description: 'Select the correct option:',
                    options: [
                        {
                            key: 'a', 
                            text: 'option A'
                        },
                        {
                            key: 'b', 
                            text: 'option B'
                        },
                        {
                            key: 'c', 
                            text: 'option C',
                            description: 'this is a sample description'
                        }
                    ]
                })


                if (selected === 'b') {
                    return 'Correct choice!'
                }

                throw 'Wrong choice!'
            },
        },
    ],
})

tool.listen()
