import { ReactNode, createContext, useContext } from "react";

interface ContextProviderProps {
  children: ReactNode;
}
export interface State {}

const RootContext = createContext<State | null>(null);

export function ContextProvider(props: ContextProviderProps) {
  return (
    <RootContext.Provider value={{}}>{props.children}</RootContext.Provider>
  );
}

export function useProvider(): State {
  const maybeContext = useContext(RootContext);
  if (!maybeContext) throw new Error("invalid application context");

  return maybeContext;
}
