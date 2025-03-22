export type TableColumnDictionary = Record<string, string>;

export type TableDisplay = {
    title?: string;
    columns: TableColumnDictionary;
    content: Array<Record<string, unknown>>;
}
