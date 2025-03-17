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
                await io.textBox(`
                    Lorem ipsum dolor sit amet. Qui quibusdam enim ut voluptatem cupiditate et unde consequatur quo dolorum nemo est porro voluptatum et dolores omnis. Non rerum dolorum et nobis ullam et perspiciatis debitis aut dolorem debitis qui voluptas veniam.
                    
                    Hic quis dolores aut alias dolorem sit impedit assumenda et minus sunt. Hic perferendis dicta ab natus enim et vero harum cum Quis neque est dolores quod. Non omnis exercitationem eos expedita itaque a molestiae eligendi ut vero deleniti et aspernatur fugiat non nihil inventore quo maxime laudantium. Ea doloribus officia et sequi molestiae At iure nobis.
                    
                    Eum corrupti nulla est ipsa consequuntur rem consequuntur consequuntur At tenetur culpa. Sed dolor minus sed facilis porro nam reiciendis fuga qui galisum ratione sed sint consectetur. Aut sapiente accusamus et eius ratione et nobis quisquam eum illo cumque. At omnis dolorum aut saepe adipisci ut voluptatibus sequi qui quis nihil sit ipsam voluptas id exercitationem voluptatem.
                `)

                return 'Everything OK!'
            },
        },
    ],
})

tool.listen()
