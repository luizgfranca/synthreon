import { useProvider } from '@/context/root'
import { ToolDto } from '@/dto/tool.dto'
import { Button } from '@/vendor/shadcn/components/ui/button'
import { Suspense } from 'react'
import { useNavigate } from 'react-router-dom'

type ProjectSidebarProps = {
    projectAcronym: string
    onSelect: (tool: ToolDto) => void
}

function ProjectSidebarContent(props: ProjectSidebarProps) {
    console.debug('ProjectSidebarContent', props);

    const navigate = useNavigate();
    const provider = useProvider();
    const tools = provider.getToolsFromProject(props.projectAcronym);

    console.debug('tools', tools);

    return (
        <div className='flex flex-1 flex-col justify-between'>
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
            <div>
                <Button 
                    variant={'outline'}
                    className='w-full flex-bottom'
                    onClick={() => navigate(`${import.meta.env.PL_PATH_PREFIX}/project/${props.projectAcronym}/tool/create/`)}
                >
                    Create New
                </Button>
            </div>
        </div>
    )
}

export function ProjectSidebar(props: ProjectSidebarProps) {
    console.debug('ProjectSidebar', props)
    // FIXME: implement a better behavior when tool query is loading
    // FIXME: handle tool query loading error
    return (
        <div className="w-full bg-zinc-900 p-2 flex flex-col flex-1">
            <div className="text-center text-zinc-100 bg-zinc-800 font-bold text-md mb-2">
                Tools
            </div>
            <nav className='flex flex-1'>
                <Suspense fallback={<h2>Loading...</h2>}>
                    <ProjectSidebarContent {...props} />
                </Suspense>
            </nav>
        </div>
    )
}
