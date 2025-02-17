import { NewToolDto, ToolDto } from "@/dto/tool.dto";
import ToolService from "@/service/tool.service";
import { Datasource } from "suspense-datasource";

export class ToolDatasource extends Datasource<ToolDto[]> {
    async fetch(projectAcronym: string): Promise<ToolDto[]> {
        return await ToolService.queryProjectTools(projectAcronym)        
    }

    async create(projectAcronym: string, data: NewToolDto) {
        // TODO: should add some logical validaions here
        await ToolService.createTool(projectAcronym, data)
        this.reset()
    }

}