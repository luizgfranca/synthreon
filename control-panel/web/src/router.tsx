import { createBrowserRouter } from "react-router-dom";
import { ProjectSelector } from "./page/projectSelector.";

const router = createBrowserRouter([
  {
    path: "/",
    element: <ProjectSelector />,
  },
]);

export { router };
