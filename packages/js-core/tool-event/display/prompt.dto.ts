export const PromptTypeOption = {
    String: "string"
};

export type PromptType = typeof PromptTypeOption[keyof typeof PromptTypeOption]

export type PromptDisplay = {
    title: string;
    type: PromptType;
};