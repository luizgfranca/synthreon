import { SelectionOption } from "./selectionOption";

type SelectionOptionParams = {
    key: string;
    title: string;
    description: string;
}

type SelectionProps = {
    options: SelectionOptionParams[];
    onSelect: (key: string) => void;
    children?: JSX.Element;
}

export default function Selection(props: SelectionProps) {
    return (
        <div className="space-y-4">
            {props.options.map((option) => (
                <SelectionOption 
                    key={option.key}
                    title={option.title}
                    description={option.description}
                    onClick={() => props.onSelect(option.key)}
                />
            ))}

            {props.children}
        </div>
    );
}
