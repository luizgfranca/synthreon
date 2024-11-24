import { ProjectDto } from "@/dto/project.dto"
import BackendService from "./backend.service"

export type QueryProjecstDto = ProjectDto[]

function queryProjects(): Promise<ProjectDto[]> {
    return new Promise((resolve) => {
        console.log('queryProjects')
        BackendService.request('/api/project')
            .then(response => response.json())
            .then(data => resolve(data as ProjectDto[]))
            .catch(e => {
                console.log(e)
                return resolve([])
            })
    })
}

const ProjectService = {
    queryProjects
}

export default ProjectService;