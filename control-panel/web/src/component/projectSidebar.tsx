import { useProvider } from '@/context/root'
import { ToolDto } from '@/dto/tool.dto'
import { Suspense } from 'react'

type ProjectSidebarProps = {
    projectAcronym: string
    onSelect: (tool: ToolDto) => void
}

function ProjectSidebarContent(props: ProjectSidebarProps) {
    console.debug('ProjectSidebarContent', props)

    const provider = useProvider()
    const tools = provider.getToolsFromProject(props.projectAcronym)

    console.debug('tools', tools)

    return (
        <ul className="space-y-1">
            {tools.map((tool) => (
                <li>
                    <div
                        className="flex items-center text-zinc-300 hover:bg-zinc-700 p-1 px-4 text-sm cursor-pointer"
                        onClick={() => props.onSelect(tool)}
                    >
                        {tool.name}
                    </div>
                </li>
            ))}
        </ul>
    )
}

export function ProjectSidebar(props: ProjectSidebarProps) {
    console.debug('ProjectSidebar', props)
    // FIXME: implement a better behavior when tool query is loading
    // FIXME: handle tool query loading error
    return (
        <div className="w-full h-full bg-zinc-900 p-2">
            <div className="text-center text-zinc-100 bg-zinc-800 font-bold text-md mb-2">
                Tools
            </div>
            <nav>
                <Suspense fallback={<h2>Loading...</h2>}>
                    <ProjectSidebarContent {...props} />
                </Suspense>
            </nav>
        </div>
    )
}
