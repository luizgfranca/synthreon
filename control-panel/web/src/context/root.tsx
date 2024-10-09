import { ProjectDatasource } from "@/datasource/project.datasource";
import { OnlyChildrenProps } from "@/lib/only-children-props";
import { ProjectDto } from "@/service/project.service";
import { createContext, useContext } from "react";

export interface State {
  getProjects: () => ProjectDto[];
}

const RootContext = createContext<State | null>(null);

const projectsDatasource = new ProjectDatasource();

export function ContextProvider(props: OnlyChildrenProps) {
  return (
    <RootContext.Provider
      value={{
        getProjects() {
          return projectsDatasource.get();
        },
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
