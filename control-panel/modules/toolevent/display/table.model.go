package tooleventdisplay

import "encoding/json"

type RowData map[string]string

type TableDisplay struct {
	Title   *string         `json:"title"`
	Columns json.RawMessage `json:"columns"`
	Content []RowData       `json:"content"`
}

func (d *TableDisplay) IsValid() bool {
	return len(d.Content) > 0
}
