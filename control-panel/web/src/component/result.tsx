import { Button } from "@/vendor/shadcn/components/ui/button";
import { CheckCircle, CircleX } from "lucide-react";


type ResultProps = {
    success: boolean;
    children: string;
    onConfirm: () => void
}

function SuccessIcon(props: { success: boolean }) {
    return props.success 
        ? <CheckCircle color="green" size={64}/>
        : <CircleX color="red" size={64}/>
}

export function Result(props: ResultProps) {
    return (
      <div>
        <div className="result grid grid-cols-5 gap-4 bg-zinc-800 p-5 rounded-md text-lg">
          <SuccessIcon success={props.success} />
          <div className="col-span-4 flex items-center">
            <div className="w-full">
              <div>{props.children}</div>
            </div>
          </div>
        </div>
        <div className="grid justify-items-end p-5 border-2 bg-zinc-900">
          <Button onClick={props.onConfirm}>OK</Button>
        </div>
      </div>
    );
    
}                       