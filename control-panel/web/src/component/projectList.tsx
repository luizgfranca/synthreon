import { useProvider } from "@/context/root";
import { OverviewSelectionButton } from "./overviewSelectionButton";
import { OverviewSelectionOption } from "./overviewSelectionOption";
import { useNavigate } from "react-router-dom";

export function ProjectList() {
  const provider = useProvider();
  const navigate = useNavigate()

  return (
    <div className="space-y-4">
      {provider.getProjects().map((project) => (
        <OverviewSelectionOption
          title={project.name}
          description={project.description}
          acronym={project.acronym}
        />
      ))}

      <OverviewSelectionButton 
        label="Create a new project" 
        onClick={() => navigate(`${import.meta.env.PL_PATH_PREFIX}/create-project`)}
      />
    </div>
  );
}
