package tooleventdisplay

type PromptType string

const (
	PromptTypeString PromptType = "string"
)

type PromptDisplay struct {
	Title string     `json:"title"`
	Type  PromptType `json:"type"`
}

func (p *PromptDisplay) IsValid() bool {
	return p.Type == PromptTypeString
}
