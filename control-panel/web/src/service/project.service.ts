import { ProjectDto } from "@/dto/project.dto"

export type QueryProjecstDto = ProjectDto[]

const BASE_URL = import.meta.env.PL_BACKEND_URL

const queryProjects: Promise<ProjectDto[]> = new Promise((resolve) => {
     fetch(`${BASE_URL}/api/project`)
        .then(response => response.json())
        .then(data => resolve(data as ProjectDto[]))
        .catch(e => {
            console.log(e)
            return resolve([])
        })
})

const ProjectService = {
    queryProjects
}

export default ProjectService;