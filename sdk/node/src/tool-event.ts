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

export type PromptType = 'string'

export type DisplayPrompt = {
    title: string;
    type: PromptType;
}

export type DisplayDefinition = {
    type: DisplayDefinitionType;
    elements?: DisplayElement[];
    result?: DisplayResult;
    prompt?: DisplayPrompt;
};

export type DsiplayRendererProps = {
    definition: DisplayDefinition
}

export type InputField = {
    name: string;
    value: string;
}

export type EventInput = {
    fields: InputField[]
}

export type EventClass = 'operation' | 'interaction' | 'announcement';
export type EventType = "display" | "input" | "open" | "provider" | "ack";


export type ToolEvent = {
    class: EventClass;
    type: EventType;
    project: string;
    tool: string;
    client?: string;
    provider?: number;
    display?: DisplayDefinition;
    input?: EventInput
}

const ToolEventUtils = {
    validate() {
        // TODO
    }
}

export { ToolEventUtils }