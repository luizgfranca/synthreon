export type TableColumnDictionary = Record<string, string>;

export type TableDisplay = {
    title?: string;
    columns: TableColumnDictionary;
    // FIXME: should support other types in the future
    content: Array<Record<string, string>>;
}
