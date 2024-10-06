import { useProvider } from "@/context/root";
import ProjectService, { ProjectDto } from "@/service/project.service";
import { OverviewSelectionOption } from "./overviewSelectionOption";
import { OverviewSelectionButton } from "./overviewSelectionButton";


const wrapPromise = (promise: Promise<unknown>) => {
    let status = "pending";
    let result: unknown;
    const suspender = promise.then(
      (r) => {
        status = "success";
        result = r;
      },
      (e) => {
        status = "error";
        result = e;
      }
    );
    return {
      read() {
        if (status === "pending") {
          throw suspender;
        } else if (status === "error") {
          throw result;
        }
        return result;
      },
    };
  };

  const projectResource = wrapPromise(ProjectService.queryProjects());

export function ProjectList() {
  const projects = projectResource.read() as ProjectDto[];


    return (
        <div className="space-y-4">
            {projects.map((project) => (
            <OverviewSelectionOption
                title={project.name}
                description={project.description}
            />
            ))}

        <OverviewSelectionButton label="Create a new project" />
      </div>
    )
}