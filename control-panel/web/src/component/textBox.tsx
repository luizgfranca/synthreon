import { Button } from "@/vendor/shadcn/components/ui/button"
import { Textarea } from "@/vendor/shadcn/components/ui/textarea"

export type TextBoxProps = {
    children: string

    onConfirm: () => void
}

export function TextBox(props: TextBoxProps) {
    return (
        <div>
            <div className="result grid bg-zinc-800 p-5 rounded-md text-lg">
                <div className="col-span-4 flex items-center font-mono">
                    <Textarea 
                        className="font-mono h-96 text-3xl" 
                        value={props.children} 
                        readOnly
                    />
                </div>
            </div>
            <div className="grid justify-items-end p-5 border-2 bg-zinc-900">
                <Button onClick={props.onConfirm}>OK</Button>
            </div>
        </div>
        
    )
}