import { RouterProvider } from "react-router-dom";
import "./index.css";
import { router } from "./router";
import { ContextProvider } from "./context/root";
import { Suspense } from "react";
import AuthService from "@/service/auth.service";

const LOGIN_PATH = `${window.location.origin}${import.meta.env.PL_PATH_PREFIX}/login`;

function App() {
    window.document.documentElement.classList.add('dark')

    if (!AuthService.isAuthenticated() && window.location.href !== LOGIN_PATH) {
        window.location.href = LOGIN_PATH;
        return <></>
    }

    return <ContextProvider>
        <Suspense fallback={<h1>Loading...</h1>}>
            <RouterProvider router={router} />
        </Suspense>
    </ContextProvider>;
}

export default App;
