import { NewToolDto, ToolDto } from "@/dto/tool.dto";
import ToolService from "@/service/tool.service";
import { Datasource } from "suspense-datasource";

export class ToolDatasource extends Datasource<ToolDto[]> {
    async fetch(projectAcronym: string): Promise<ToolDto[]> {
        // FIXME: temporary, disabling the cache until better implementation
        // that will invalidate the cache when:
        //  - a certain time has passed
        //  - there was a timeout
        this.reset();
        return await ToolService.queryProjectTools(projectAcronym);
    }

    async create(projectAcronym: string, data: NewToolDto) {
        // TODO: should add some logical validaions here
        await ToolService.createTool(projectAcronym, data);
        this.reset();
    }

}
