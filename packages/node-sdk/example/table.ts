import { ToolProvider } from '../src/tool-provider'

const columns = {
    food: 'Food',
    annualQty: 'Annual Quantities',
    annualCost: 'Annual Cost',
}

const data = [
    {
        food: 'Wheat Flour',
        annualQty: '370 lb (170 kg)',
        annualCost: '$13.33',
    },
    {
        food: 'Evaporated Milk',
        annualQty: '57 cans',
        annualCost: '$3.84',
    },
    {
        food: 'Cabbag',
        annualQty: '111 lb (50 kg)',
        annualCost: '$4.11',
    },
    {
        food: 'Spinach',
        annualQty: '23 lb (10 kg)',
        annualCost: '$1.85',
    },
    {
        food: 'Dried Navy Beans',
        annualQty: '285 lb (129 kg)',
        annualCost: '$16.80',
    }
]


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
                await io.table({
                    title: "Stigler's 1939 Diet",
                    content: data,
                    columns: columns
                }) 

                return 'done!'
            },
        },
    ],
})

tool.listen()

