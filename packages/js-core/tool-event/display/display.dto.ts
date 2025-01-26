import { InformationDisplay } from "./information.dto";
import { PromptDisplay } from "./prompt.dto";
import { TextBoxDisplay } from "./textbox.dto";

export const DisplayTypeValue = {
    Prompt: "prompt",
    Information: "information",
    TextBox: "textbox"
} as const;

export type DisplayType = typeof DisplayTypeValue[keyof typeof DisplayTypeValue]

export type DisplayDefinition = {
    type: DisplayType;
    prompt?: PromptDisplay;
    information?: InformationDisplay;
    textBox?: TextBoxDisplay;
};

export * from './prompt.dto'
export * from './information.dto'
export * from './textbox.dto'