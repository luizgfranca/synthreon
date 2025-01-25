import { InformationDisplay } from "./information.dto";
import { PromptDisplay } from "./prompt.dto";
import { TextBoxDisplay } from "./textbox.dto";

const DisplayTypeValue = {
    Prompt: "prompt",
    Information: "information",
    TextBox: "textbox"
} as const;

type DisplayType = typeof DisplayTypeValue[keyof typeof DisplayTypeValue]

export type DisplayDefinition = {
    type: DisplayType;
    prompt?: PromptDisplay;
    information?: InformationDisplay;
    textBox?: TextBoxDisplay;
};