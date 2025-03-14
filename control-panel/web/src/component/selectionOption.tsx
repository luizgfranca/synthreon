export interface SelectionOptionProps {
    title: string;
    description: string;
    onClick: () => void;
}

export function SelectionOption(props: SelectionOptionProps) {
    return (
        <div
            className="bg-zinc-800 hover:bg-zinc-700 p-4 cursor-pointer"
            onClick={props.onClick}
        >
            <h2 className="text-xl font-semibold mb-2">{props.title}</h2>
            <p className="text-zinc-400">{props.description}</p>
        </div>
    );
}
