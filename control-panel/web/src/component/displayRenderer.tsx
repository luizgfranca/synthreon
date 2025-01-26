import { useCallback, useMemo, useState } from "react";
import { Result } from "./result";
import { Prompt } from "./prompt";
import { TextBox } from "./textBox";
import { ToolEventResult } from "platformlab-core/dist/tool-event/result/result.dto";
import { ToolEventDto } from "platformlab-core";

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

type DisplayPrompt = {
    title: string;
    type: string;
}

type DisplayTextBox = {
    content: string
}

type DisplayDefinitionType = 'result' | 'view' | 'prompt' | 'textbox' | string

export type DisplayDefinition = {
    type: DisplayDefinitionType;
    elements?: DisplayElement[];
    result?: DisplayResult;
    prompt?: DisplayPrompt;
    textBox?: DisplayTextBox;
};

export type Field = {
    name: string,
    value: string
}

export type DsiplayRendererProps = {
    event: ToolEventDto
    onSumission: (fields: Field[]) => void
    
    resetCallback: () => void
}

export function DisplayRenderer(props: DsiplayRendererProps) {
    console.debug('on displayRenderer', props.event)

    if (props.event.type === 'command/finish') {
        const success = props.event.result && props.event.result.status === 'success' || false

        return (
            <Result 
                success={success}
                onConfirm={() => props.resetCallback()}
            >
                {props.event.result?.message ?? ''}
            </Result>
        )
    }

    if(!props.event.display) {
        throw new Error('expected display but its not defined')
    }

    switch(props.event?.display.type) {
        case 'prompt':
            return (
                <Prompt 
                    title={props.event.display.prompt?.title ?? ''} 
                    onSubmit={(value) => props.onSumission([{
                        name: 'prompt',
                        value
                    }])}
                />
            )
        case 'textbox':
            return (
                <TextBox
                    onConfirm={() => props.onSumission([])}
                >
                    {props.event.display.textBox?.content ?? ''}
                </TextBox>
            )
    }

}