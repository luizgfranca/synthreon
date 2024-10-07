import { wrapPromise } from "@/lib/suspense-wrapper";
import ProjectService from "@/service/project.service";

const projectResource = wrapPromise(ProjectService.queryProjects);

export {
    projectResource
}