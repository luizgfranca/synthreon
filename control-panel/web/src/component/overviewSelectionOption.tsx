import { useNavigate } from "react-router-dom";

export interface OverviewSelectionOptionProps {
  title: string;
  description: string;
  acronym: string;
}

export function OverviewSelectionOption(props: OverviewSelectionOptionProps) {
  const navigate = useNavigate();

  return (
    <div
      className="bg-zinc-800 hover:bg-zinc-700 p-4 cursor-pointer"
      onClick={() => navigate(`/project/${props.acronym}`)}
    >
      <h2 className="text-xl font-semibold mb-2">{props.title}</h2>
      <p className="text-zinc-400">{props.description}</p>
    </div>
  );
}
