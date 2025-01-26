export const InformationDisplayType = {
    Success: "success",
    Failure: "failure"
};

type InformationDisplayType = typeof InformationDisplayType[keyof typeof InformationDisplayType]

export type InformationDisplay = {
    type: InformationDisplayType;
    message: string;
}