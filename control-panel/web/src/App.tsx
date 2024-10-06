import { RouterProvider } from "react-router-dom";
import "./index.css";
import { router } from "./router";
import { ContextProvider } from "./context/root";
import { Suspense } from "react";

function App() {
  return <ContextProvider>
    <Suspense fallback={<h1>Loading...</h1>}>
      <RouterProvider router={router} />
    </Suspense>
  </ContextProvider>;
}

export default App;
