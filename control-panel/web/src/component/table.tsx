export type TableColumnDictionary = Record<string, string>;

export type TableProps = {
    title?: string;
    columns: TableColumnDictionary;
    content: Array<Record<string, unknown>>;
}

function MaybeTitle(props: {title?: string}) {
    return props.title 
        ? (<h2 className="text-center text-xl mb-4">{props.title}</h2>) 
        : null
}


export default function Table(props: TableProps) {
    console.debug('rendering table', props);
    return (
        <div className="w-full">
            <MaybeTitle title={props.title} />
            <table className="w-full border-collapse bg-zinc-1000 shadow-lg font-mono">
                <thead className='border border-zinc-700'>
                    <tr className="bg-zinc-900 border-b-2 border-zinc-1000">
                        {
                            Object.keys(props.columns)
                                .map(key => (
                                    <th className="py-3 px-4 text-left text-zinc-50 font-bold uppercase ">{props.columns[key]}</th>
                                ))
                        }
                    </tr>
                </thead>
                <tbody className="border border-zinc-700">
                    {
                        props.content.map((item, i) => (
                            <tr className={"border-b border-zinc-700 " + (i % 2 !== 0 ? 'bg-zinc-900/50' : '')}>
                                {
                                    Object.keys(props.columns).map(key => (
                                        <td className="py-2 px-4 border-r border-zinc-800 text-zinc-300">{item[key] as string}</td>
                                    ))
                                }
                            </tr>
                        ))
                    }
                </tbody>
            </table>
        </div>
    )
}
