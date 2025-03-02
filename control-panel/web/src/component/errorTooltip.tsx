import { AlertCircle } from "lucide-react"

import {
    Alert,
    AlertDescription,
    AlertTitle,
} from "@/vendor/shadcn/components/ui/alert"

type ErrorTooltipProps = {
    message: string;
}


export function ErrorTooltip(props: ErrorTooltipProps) {
    return (
        <Alert variant="destructive" className="text-red-600">
            <AlertCircle stroke="red" className="h-4 w-4 border-blue-100" />
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>
                {props.message}
            </AlertDescription>
        </Alert>
    )
}
