import { Button } from "@/vendor/shadcn/components/ui/button";

export type DisplayBoxProps = {
    children: string | JSX.Element;
    onConfirm: () => void
}


export default function DisplayBox(props: DisplayBoxProps) {
    return (
        <div className="border border-zinc-700 shadow-lg">
            <div className="result rounded-md text-lg border-zinc-900">
                <div className="p-5 flex items-center">
                    {props.children} 
                </div>
             <div className="grid justify-items-end p-5 border-t-0 bg-zinc-900/90">
                <Button className="rounded-none px-8" onClick={props.onConfirm}>OK</Button>
            </div>
        </div>
    </div>
    )
}
