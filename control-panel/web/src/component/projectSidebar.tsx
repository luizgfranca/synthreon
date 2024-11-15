import { ToolDto } from "@/dto/tool.dto";

type ProjectSidebarProps = {
  tools: ToolDto[];
  onSelect: (tool: ToolDto) => void;
};

export function ProjectSidebar(props: ProjectSidebarProps) {
  return (
    <div className="w-full h-full bg-zinc-900 p-2">
      <div className="text-center text-zinc-100 bg-zinc-800 font-bold text-md mb-2">
        Tools
      </div>
      <nav>
        <ul className="space-y-1">
          {props.tools.map((tool) => (
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
      </nav>
    </div>
  );
}
