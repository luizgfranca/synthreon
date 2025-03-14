import { useProvider } from "@/context/root";
import { OverviewSelectionButton } from "./overviewSelectionButton";
import { useNavigate } from "react-router-dom";
import Selection from "./selection";

export function ProjectList() {
    const provider = useProvider();
    const navigate = useNavigate()
    // TODO: create navigator that already computes prefix path

    return (
        <Selection
            options={
                provider.getProjects().map((project) => ({
                key: project.acronym,
                title: project.name,
                description: project.description
            }))}
            onSelect={(key) => navigate(`${import.meta.env.PL_PATH_PREFIX}/project/${key}`)}
        >
            <OverviewSelectionButton
                label="Create a new project"
                onClick={() => navigate(`${import.meta.env.PL_PATH_PREFIX}/create-project`)}
            />
        </Selection>
    )
}
