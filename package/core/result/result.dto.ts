export const ToolEventResultStatus = {
    Success: "success",
    Failure: "failure",
} as const;
  
export type ToolEventResultStatus = typeof ToolEventResultStatus[keyof typeof ToolEventResultStatus];

export type ToolEventResult = {
    status: ToolEventResultStatus;
    message: string;
};