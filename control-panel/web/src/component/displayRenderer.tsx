import { ToolEventDto } from "platformlab-core";
import { Prompt } from "./prompt";
import { Result } from "./result";
import { TextBox } from "./textBox";

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