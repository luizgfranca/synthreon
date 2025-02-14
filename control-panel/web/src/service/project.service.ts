import { NewProjectDto, ProjectDto } from "@/dto/project.dto"
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

function createProject(dto: NewProjectDto) {
    return new Promise<void>((resolve, reject) => {
        // TODO restructure the backendservice to better
        // abstract many types of requests
        BackendService.request(
            '/api/project',
            {
                method: 'POST',
                body: JSON.stringify(dto)
            }
        )
        .then(() => resolve())
        .catch(e => {
            console.log(e)
            return reject(e)
        })
    })
}

const ProjectService = {
    queryProjects,
    createProject
}

export default ProjectService;