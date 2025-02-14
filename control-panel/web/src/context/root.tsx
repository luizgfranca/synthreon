import { ProjectDatasource } from "@/datasource/project.datasource";
import { NewProjectDto, ProjectDto } from "@/dto/project.dto";
import { OnlyChildrenProps } from "@/lib/only-children-props";
import { createContext, useContext } from "react";

// TODO: refactor this to expose datasource directly
export interface State {
  getProjects: () => ProjectDto[];
  createProject: (data: NewProjectDto) => Promise<void>
}

const RootContext = createContext<State | null>(null);

const projectsDatasource = new ProjectDatasource();

function getProjects() {
  return projectsDatasource.get();
}

async function createProject(data: NewProjectDto) {
  return await projectsDatasource.create(data);
}

export function ContextProvider(props: OnlyChildrenProps) {
  return (
    <RootContext.Provider
      value={{
        getProjects,
        createProject
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
