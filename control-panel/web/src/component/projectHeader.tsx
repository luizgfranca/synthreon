import { useNavigate } from "react-router-dom";

type ProjectHeaderProps = {
    projectName: string;
    onLogoutClick: () => void
}

export function ProjectHeader(props: ProjectHeaderProps) {
    const navigate = useNavigate();

    const navigateHome = () => navigate(`${import.meta.env.PL_PATH_PREFIX}`);

    return (
        <div className="flex items-center justify-between bg-zinc-900 h-12 px-2 w-100 border">
            <div className="flex items-center cursor-pointer" onClick={navigateHome}>
                <span className="text-lg font-semibold">Synthreon</span>
            </div>

            <button 
                className="mr-10 px-5 py-2 bg-zinc-800 text-zinc-100 font-bold" 
                onClick={navigateHome}
            >
                <span>{props.projectName}</span>
            </button>

            <div className="flex items-center">
                <button className="text-sm mr-10 px-2 py-1 bg-zinc-800 text-zinc-100" onClick={props.onLogoutClick}>
                    <span>Logout</span>
                </button>
            </div>
        </div>
    );
}
