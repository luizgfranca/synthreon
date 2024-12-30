package tooleventdisplay

type TextBoxDisplay struct {
	Content string `json:"content"`
}

func (t *TextBoxDisplay) IsValid() bool {
	return true
}
