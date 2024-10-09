import ProjectService, { ProjectDto } from "@/service/project.service";

type LoadingStatus = 'none' | 'pending' | 'success' | 'fail';

export class ProjectDatasource {
    #loadingStatus: LoadingStatus = 'none';
    #loadingPromise?: Promise<ProjectDto[]>;
    
    #data?: ProjectDto[];
    #error: unknown;
    
    async fetch(): Promise<ProjectDto[]>{
        return ProjectService.queryProjects;
    }

    get(): ProjectDto[] {
        switch(this.#loadingStatus) {
            case 'none':
                this.#loadingPromise = this.fetch();
                this.#loadingPromise.
                    then((data) => {
                        this.#data = data;
                        this.#loadingStatus = 'success'
                    }).catch(e => {
                        this.#error = e;
                        this.#loadingStatus = 'fail'
                    });
                this.#loadingStatus = 'pending'
                throw this.#loadingPromise;
            case 'pending':
                throw this.#loadingPromise;
            case 'fail':
                throw this.#error
            case 'success':
                return this.#data as ProjectDto[];
        }
    }
}