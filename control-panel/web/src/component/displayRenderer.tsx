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

export type DisplayDefinition = {
    type: DisplayDefinitionType;
    elements?: DisplayElement[];
    result?: DisplayResult;
};

export type DsiplayRendererProps = {
    definition: DisplayDefinition
}

export function DisplayRenderer(props: DsiplayRendererProps) {
    return (
        <Result success={props.definition.result?.success ?? false}>{props.definition.result?.message ?? ''}</Result>
    )
}