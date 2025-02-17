import { NewToolDto, ToolDto } from "@/dto/tool.dto";
import BackendService from "./backend.service"

async function queryProjectTools(projectAcronym: string): Promise<ToolDto[]> {
    const result = await BackendService.request(`/api/project/${projectAcronym}/tool`)
    return await result.json() as ToolDto[]
}

async function createTool(projectAcronym: string, data: NewToolDto) {
    // TODO restructure the backendservice to better
    //      abstract many types of requests
    try {
        await BackendService.request(`/api/project/${projectAcronym}/tool`, {
            method: 'POST',
            body: JSON.stringify(data)
        })
    } catch(e) {
        console.log(e)
        throw e
    }
}

const ToolService = {
    queryProjectTools,
    createTool
}

export default ToolService;