import { createBrowserRouter } from "react-router-dom";
import { LoginPage } from "./page/login.page";
import { ProjectOverviewPage } from "./page/projectOverview.page";
import { ProjectSelector } from "./page/projectSelector.page";
import CreateProjectPage from "./page/createProject.page";

const router = createBrowserRouter([
  {
    path: `${import.meta.env.PL_PATH_PREFIX}`,
    element: <ProjectSelector />,
  },
  {
    path: `${import.meta.env.PL_PATH_PREFIX}/login`,
    element: <LoginPage />,
  },
  {
    path: `${import.meta.env.PL_PATH_PREFIX}/project/:projectAcronym`,
    element: <ProjectOverviewPage />,
  },
  {
    path: `${import.meta.env.PL_PATH_PREFIX}/create-project/`,
    element: <CreateProjectPage />,
  }
]);

export { router };
