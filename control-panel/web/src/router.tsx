import { createBrowserRouter } from "react-router-dom";
import { LoginPage } from "./page/login.page";
import { ProjectOverviewPage } from "./page/projectOverview.page";
import { ProjectSelector } from "./page/projectSelector.page";

const router = createBrowserRouter([
  {
    path: "/",
    element: <ProjectSelector />,
  },
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    path: '/project/*',
    element: <ProjectOverviewPage />,
  }
]);

export { router };
