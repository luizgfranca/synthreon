export type SelectionOption = {
    key: string;
    text: string;
    description?: string;
}

export type SelectionDisplay  = {
    description: string;
    options: SelectionOption[]
}

