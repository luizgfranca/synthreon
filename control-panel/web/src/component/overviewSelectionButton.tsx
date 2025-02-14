interface OverviewSelectionButtonProps {
    label: string,
    onClick: () => void
}

export function OverviewSelectionButton(props: OverviewSelectionButtonProps) {
    return (
        <div 
            className="bg-zinc-900 hover:bg-zinc-700 p-4 text-center cursor-pointer"
            onClick={props.onClick}
        >
            <i className="fas fa-plus mr-2 text-white"></i>{props.label}
        </div>
    )    
}