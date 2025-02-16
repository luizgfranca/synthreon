import { ToolDto } from "@/dto/tool.dto";
import ToolService from "@/service/tool.service";
import { Datasource } from "suspense-datasource";

export class ToolDatasource extends Datasource<ToolDto[]> {
    async fetch(projectAcronym: string): Promise<ToolDto[]> {
        console.debug('ToolDatasource fetch', projectAcronym)
        const result = await ToolService.queryProjectTools(projectAcronym)        
        console.debug('ToolDatasource fetch result', result)
        return result
    }
}