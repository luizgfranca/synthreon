import { SelectionOption } from "./selectionOption";

type SelectionOptionParams = {
    key: string;
    title: string;
    description?: string;
}

type SelectionProps = {
    introduction?: string;
    options: SelectionOptionParams[];
    onSelect: (key: string) => void;
    children?: JSX.Element;
}

function maybeIntroduction(introductionParam?: string) {
    return introductionParam ? (
        <div className="space-y-1">
            {introductionParam}
        </div>
    ) : <></>
}

export default function Selection(props: SelectionProps) {
    return (
        <div className="space-y-4">
            {maybeIntroduction(props.introduction)}
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
