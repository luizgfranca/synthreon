import { useProvider } from '@/context/root'
import { ToolDto } from '@/dto/tool.dto'
import { Button } from '@/vendor/shadcn/components/ui/button'
import { Suspense, useCallback, useState } from 'react'
import { useNavigate } from 'react-router-dom'

type ProjectSidebarProps = {
    projectAcronym: string
    onSelect: (tool: ToolDto) => void
}

const TOOL_ITEM_STYLE_BASE = "flex items-center text-zinc-300 hover:bg-zinc-700 p-1 px-4 text-sm cursor-pointer";
const TOOL_ITEM_STYLE_SELECTED = "bg-zinc-800";

function getItemStyle(selected: boolean) {
    return selected
        ? `${TOOL_ITEM_STYLE_BASE} ${TOOL_ITEM_STYLE_SELECTED}`
        : TOOL_ITEM_STYLE_BASE
}


function ProjectSidebarContent(props: ProjectSidebarProps) {
    console.debug('ProjectSidebarContent', props);

    const navigate = useNavigate();
    const provider = useProvider();
    const tools = provider.getToolsFromProject(props.projectAcronym);
    const [selectedToolAcronym, setSeletectedToolAcronym] = useState<string | null>(null)

    const selectTool = useCallback((tool: ToolDto) => {
        console.debug('toolSelected', tool);
        setSeletectedToolAcronym(tool.acronym);
        props.onSelect(tool)
    }, [props, setSeletectedToolAcronym])

    console.debug('tools', tools);
    console.debug('selected', selectedToolAcronym);

    return (
        <div className='flex flex-1 flex-col justify-between'>
            <ul className="space-y-1">
                {tools.map((tool) => (
                    <li>
                        <div
                            key={tool.acronym}
                            className={getItemStyle(tool.acronym === selectedToolAcronym)}
                            onClick={() => selectTool(tool)}
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
