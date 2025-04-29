import { NewToolDto, ToolDto } from "@/dto/tool.dto";
import ToolService from "@/service/tool.service";
import { Datasource } from "suspense-datasource";

export class ToolDatasource extends Datasource<ToolDto[]> {
    async fetch(projectAcronym: string): Promise<ToolDto[]> {
        return await ToolService.queryProjectTools(projectAcronym)
            .finally(() => {
                // FIXME: temporary, disabling the cache until better implementation
                // that will invalidate the cache when:
                //  - a certain time has passed
                //  - there was a timeout
                //  needs the setTimeout to put the reset on the back of the microtask
                //  queue instead of running it before the valua can be effectivelly
                //  returned by the get
                setTimeout(() => this.reset(projectAcronym), 1);
            });
    }

    async create(projectAcronym: string, data: NewToolDto) {
        // TODO: should add some logical validaions here
        await ToolService.createTool(projectAcronym, data);
        this.reset();
    }

}
