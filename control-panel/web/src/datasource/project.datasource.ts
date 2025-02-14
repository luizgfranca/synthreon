import { NewProjectDto, ProjectDto } from "@/dto/project.dto";
import { Datasource } from 'suspense-datasource';
import ProjectService from "@/service/project.service";

export class ProjectDatasource extends Datasource<ProjectDto[]> {
    async fetch(): Promise<ProjectDto[]>{
        return ProjectService.queryProjects();
    }

    async create(data: NewProjectDto) {
        // FIXME: add validations for the project here
        await ProjectService.createProject(data)
        this.reset()
    }
}