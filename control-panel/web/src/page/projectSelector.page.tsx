import { ProjectList } from "@/component/projectList";
import AuthService from "@/service/auth.service";
import { Suspense } from "react";
import { useNavigate } from "react-router-dom";

export function ProjectSelector() {
  const navigate = useNavigate();

  if(!AuthService.isAuthenticated()) {
    navigate('/login')
  }

  return (
    <div className="bg-zinc-900 text-zinc-100 h-screen">
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-6">Projects</h1>

        <Suspense
          fallback={<h2 className="text-xl font-bold mb-6">Loading...</h2>}
        >
          <div className="space-y-4">
            <ProjectList />
          </div>
        </Suspense>
      </div>
    </div>
  );
}
