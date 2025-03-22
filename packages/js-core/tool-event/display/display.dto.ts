import { InformationDisplay } from "./information.dto";
import { PromptDisplay } from "./prompt.dto";
import { SelectionDisplay } from "./selection.dto";
import { TableDisplay } from "./table.dto";
import { TextBoxDisplay } from "./textbox.dto";

export const DisplayTypeValue = {
    Prompt: "prompt",
    Information: "information",
    TextBox: "textbox",
    Selection: "selection",
    Table: "table"
} as const;

export type DisplayType = typeof DisplayTypeValue[keyof typeof DisplayTypeValue]

export type DisplayDefinition = {
    type: DisplayType;
    prompt?: PromptDisplay;
    information?: InformationDisplay;
    textBox?: TextBoxDisplay;
    selection?: SelectionDisplay;
    table?: TableDisplay;
};

export * from './prompt.dto'
export * from './information.dto'
export * from './textbox.dto'
export * from './selection.dto'
export * from './table.dto'
