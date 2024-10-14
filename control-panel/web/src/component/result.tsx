import { CheckCircle } from "lucide-react";


type ResultProps = {
    success: boolean;
    children: string;
}

export function Result(props: ResultProps) {
    return (
        <div className="result grid grid-cols-5 gap-4 bg-zinc-800 p-5 rounded-md text-lg">
            <CheckCircle color="green" size={64}/>
            <div className="col-span-4 flex items-center">
                <span>{props.children}</span>
            </div>
        </div>              
    )
    
}                       