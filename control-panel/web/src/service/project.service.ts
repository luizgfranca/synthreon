export type ProjectDto = {
    id: number,
    acronym: string,
    name: string,
    description: string
}

export type QueryProjecstDto = ProjectDto[]

const BASE_URL = import.meta.env.PL_BACKEND_URL

const queryProjects = new Promise((resolve) => {
     fetch(`${BASE_URL}/project`)
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