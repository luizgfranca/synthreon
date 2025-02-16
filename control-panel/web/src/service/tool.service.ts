import { ToolDto } from "@/dto/tool.dto";
import BackendService from "./backend.service"

async function queryProjectTools(projectAcronym: string): Promise<ToolDto[]> {
    const result = await BackendService.request(`/api/project/${projectAcronym}/tool`)
    return await result.json() as ToolDto[]
}

const ToolService = {
    queryProjectTools,
}

export default ToolService;