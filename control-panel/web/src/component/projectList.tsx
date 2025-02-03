import { useProvider } from "@/context/root";
import { OverviewSelectionButton } from "./overviewSelectionButton";
import { OverviewSelectionOption } from "./overviewSelectionOption";

export function ProjectList() {
  const provider = useProvider();

  return (
    <div className="space-y-4">
      {provider.getProjects().map((project) => (
        <OverviewSelectionOption
          title={project.name}
          description={project.description}
          acronym={project.acronym}
        />
      ))}

      <OverviewSelectionButton label="Create a new project" />
    </div>
  );
}
