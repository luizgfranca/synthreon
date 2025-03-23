package tooleventdisplay

type TableColumnDictionary map[string]string
type RowData map[string]string

// FIXME: table column order is not deterministic, when reaching the go side
//        it can change the order because it does not garantee the ordering of maps
type TableDisplay struct {
	Title   *string                `json:"title"`
	Columns *TableColumnDictionary `json:"columns"`
	Content []RowData              `json:"content"`
}

func (d *TableDisplay) IsValid() bool {
	return len(d.Content) > 0
}
