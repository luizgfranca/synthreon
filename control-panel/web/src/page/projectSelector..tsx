import { OverviewSelectionButton } from "@/component/overviewSelectionButton";
import { OverviewSelectionOption } from "@/component/overviewSelectionOption";

const projects = ['Project A', 'Project B', 'Project C']

export function ProjectSelector() {
    return (
        <div className="bg-zinc-900 text-zinc-100 h-screen">
          <div className="container mx-auto px-4 py-8">
            <h1 className="text-3xl font-bold mb-6">Projects</h1>
            
            <div className="space-y-4">
              {
                projects.map(project => <OverviewSelectionOption title={project} description="this is the project description"/>)
              }  
    
              <OverviewSelectionButton label="Create a new project" />
            </div>
          </div>
        </div>
      );
}