import { ToolEventDto } from "@synthreon/core";
import { Prompt } from "./prompt";
import { Result } from "./result";
import { TextBox } from "./textBox";
import Selection from "./selection";
import DisplayBox from "./displayBox";
import Table from "./table";

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

    if (!props.event.display) {
        throw new Error('expected display but its not defined')
    }

    let definition;
    switch (props.event?.display.type) {
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
        case "selection":
            definition = props.event.display?.selection;

            if (!definition) {
                throw new Error('select specified but no definition for it found')
            }
            return (
                <Selection
                    introduction={definition.description}
                    options={definition.options.map((option) => ({
                        key: option.key,
                        title: option.text,
                        description: option.description
                    }))}
                    onSelect={(key) => props.onSumission([{
                        name: 'selection',
                        value: key
                    }])}
                />
            )
        case "table":
            definition = props.event.display?.table;
            if (!definition) {
                throw new Error('table specified but no definition for it found')
            }

            return (
                <DisplayBox
                    onConfirm={() => props.onSumission([])}
                >
                    <Table {...definition} />
                </DisplayBox>
            )
    }

}
