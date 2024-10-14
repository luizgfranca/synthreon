import { useCallback, useMemo, useState } from "react";
import { Result } from "./result";

type DisplayElement = {
    type: string;
    label: string;
    text: string;
    description: string;
    name: string;
};

type DisplayResult = {
    success: boolean;
    message: string;
};

type DisplayDefinitionType = 'result' | 'view' | 'prompt' | string

type DisplayDefinition = {
    type: DisplayDefinitionType;
    elements?: DisplayElement[];
    result?: DisplayResult;
};

type DsiplayRendererProps = {
    definition: DisplayDefinition
}

export function DisplayRenderer(props: DsiplayRendererProps) {
    const socket = useMemo(() => {
        const s = new WebSocket(`${import.meta.env.PL_BACKEND_URL}/api/tool/ws`)
        s.addEventListener('open', () => {
            console.log('socket open')
            setInterval(() => {
                s.send(JSON.stringify({
                    "class": "interaction",
                    "type": "open",
                    "project": "proj-x",
                    "tool": "tool-y",
                }))
                console.log('sent')
            }, 1000)
        })

        s.addEventListener('message', (e) => {
            console.log(`recv: ${e.data}`)
        })

        return s    
    }, [])

    return (
        <Result success={props.definition.result?.success ?? false}>{props.definition.result?.message ?? ''}</Result>
    )
}