import { InformationDisplay } from "./information.dto";
import { PromptDisplay } from "./prompt.dto";
import { TextBoxDisplay } from "./textbox.dto";
declare const DisplayTypeValue: {
    readonly Prompt: "prompt";
    readonly Information: "information";
    readonly TextBox: "textbox";
};
type DisplayType = typeof DisplayTypeValue[keyof typeof DisplayTypeValue];
export type DisplayDefinition = {
    type: DisplayType;
    prompt?: PromptDisplay;
    information?: InformationDisplay;
    textBox?: TextBoxDisplay;
};
export {};
