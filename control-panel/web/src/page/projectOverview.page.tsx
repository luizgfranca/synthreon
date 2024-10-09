export function ProjectOverviewPage(props: unknown) {
    console.log('props', props)
    const project = window.location.pathname.split('/')[2];


    return (
        <div className="bg-zinc-900 text-zinc-100 h-screen">
            <div className="container mx-auto px-4 py-8">
                <h1 className="text-3xl font-bold mb-6">{project}</h1>

            </div>
        </div>
    )
}