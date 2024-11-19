import { createBrowserRouter } from "react-router-dom";
import { ProjectSelector } from "./page/projectSelector.page";
import { ProjectOverviewPage } from "./page/projectOverview.page";
import { ToolViewExperimentsPage } from "./page/toolViewExperiments.page";
import { LoginPage } from "./page/login.page";

const router = createBrowserRouter([
  {
    path: "/",
    element: <ProjectSelector />,
  },
  {
    path: "/auth/login",
    element: <LoginPage />,
  },
  {
    path: '/project/*',
    element: <ProjectOverviewPage />,
  },
  {
    path: '/exp',
    element: <ToolViewExperimentsPage />
  }
]);

export { router };
