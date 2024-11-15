type EmptyStateViewProps = {
    children: string
}

export function EmptyState(props: EmptyStateViewProps) {
    return (
        <div className="text-zinc-100 w-full">
            <div className="container mx-auto px-4 py-8">
                <h1 className="text-xl text-center font-light mb-6">{props.children}</h1>
            </div>
        </div>
    )
}