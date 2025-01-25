export const PromptType = {
    String: "string"
};

type PromptType = typeof PromptType[keyof typeof PromptType]

export type PromptDisplay = {
    title: string;
    type: PromptType;
};