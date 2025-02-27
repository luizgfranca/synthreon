import { ProjectDatasource } from "@/datasource/project.datasource";
import { ToolDatasource } from "@/datasource/tool.datasource";
import { NewProjectDto, ProjectDto } from "@/dto/project.dto";
import { NewToolDto, ToolDto } from "@/dto/tool.dto";
import { OnlyChildrenProps } from "@/lib/only-children-props";
import { createContext, useContext } from "react";

// TODO: refactor this to expose datasource directly
export interface State {
  getProjects: () => ProjectDto[];
  createProject: (data: NewProjectDto) => Promise<void>,
  getToolsFromProject: (projectAcronym: string) => ToolDto[]
  createTool: (projectAcronym: string, data: NewToolDto) => Promise<void>
}

const RootContext = createContext<State | null>(null);

const projectsDatasource = new ProjectDatasource();
const toolDatasource = new ToolDatasource();

function getProjects() {
    return projectsDatasource.get();
}

function getToolsFromProject(projectAcronym: string): ToolDto[] {
    return toolDatasource.get(projectAcronym);
}

async function createProject(data: NewProjectDto) {
    return await projectsDatasource.create(data);
}

async function createTool(projectAcronym: string, data: NewToolDto) {
    return await toolDatasource.create(projectAcronym, data);
}

export function ContextProvider(props: OnlyChildrenProps) {
    return (
        <RootContext.Provider
            value={{
                getProjects,
                createProject,
                getToolsFromProject,
                createTool
            }}
        >
        {props.children}
        </RootContext.Provider>
  );
}

export function useProvider(): State {
  const maybeContext = useContext(RootContext);
  if (!maybeContext) throw new Error("invalid application context");

  return maybeContext;
}
