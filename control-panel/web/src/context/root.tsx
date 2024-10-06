import ProjectService, { ProjectDto } from "@/service/project.service";
import { ReactNode, createContext, useContext, useState } from "react";

interface ContextProviderProps {
  children: ReactNode;
}
export interface State {
  // getProjects: () => ProjectDto[]
}

type Query<T> = {
  status: 'pending' | 'success' | 'error',
  data?: T | Promise<T>
}

const RootContext = createContext<State | null>(null);

export function ContextProvider(props: ContextProviderProps) {
  const [projects, setProjects] = useState<Query<ProjectDto[]>>()

  // const getProjects = (): ProjectDto[] => {
  //   if (!projects) {
  //     const promise = ProjectService.queryProjects()
  //     setProjects({status: 'pending', data: promise});
      
  //     promise.then(projects => {
  //       console.debug('setProjects', {status: 'success', data: projects})
  //       setProjects({status: 'success', data: projects})
  //       return projects
  //     })
  //     .catch(() => {
  //       setProjects({status: 'error', data: []})
  //       return []
  //     });

  //     throw promise
  //   }

    // switch (projects.status) {
    //   case 'pending':
    //   case 'error':
    //     console.debug('pending/error')
    //     throw projects.data
      
    //   case 'success':
    //     console.debug('success')
    //     return projects.data as ProjectDto[]
    // }
  // }

  return (
    <RootContext.Provider value={{
      
    }}>{props.children}</RootContext.Provider>
  );
}

export function useProvider(): State {
  const maybeContext = useContext(RootContext);
  if (!maybeContext) throw new Error("invalid application context");

  return maybeContext;
}
