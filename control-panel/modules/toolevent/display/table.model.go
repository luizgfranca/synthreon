package tooleventdisplay

type TableColumnDictionary map[string]string
type RowData map[string]string

type TableDisplay struct {
	Title   *string                `json:"title"`
	Columns *TableColumnDictionary `json:"columns"`
	Content []RowData              `json:"content"`
}

func (d *TableDisplay) IsValid() bool {
	return len(d.Content) > 0
}
