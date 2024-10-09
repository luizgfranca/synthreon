import { createBrowserRouter } from "react-router-dom";
import { ProjectSelector } from "./page/projectSelector.page";
import { ProjectOverviewPage } from "./page/projectOverview.page";

const router = createBrowserRouter([
  {
    path: "/",
    element: <ProjectSelector />,
  },
  {
    path: '/project/*',
    element: <ProjectOverviewPage />,
  }
]);

export { router };
