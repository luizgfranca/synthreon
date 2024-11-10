import { CheckCircle, CircleX } from "lucide-react";


type ResultProps = {
    success: boolean;
    children: string;
}

function SuccessIcon(props: { success: boolean }) {
    return props.success 
        ? <CheckCircle color="green" size={64}/>
        : <CircleX color="red" size={64}/>
}

export function Result(props: ResultProps) {
    return (
        <div className="result grid grid-cols-5 gap-4 bg-zinc-800 p-5 rounded-md text-lg">
            <SuccessIcon success={props.success} />
            <div className="col-span-4 flex items-center">
                <span>{props.children}</span>
            </div>
        </div>              
    )
    
}                       