package tooleventdisplay

type SelectionOption struct {
	Key         string `json:"key"`
	Text        string `json:"text"`
	Description string `json:"description"`
}

type SelectionDisplay struct {
	Description string            `json:"description"`
	Options     []SelectionOption `json:"options"`
}

func (o *SelectionOption) IsValid() bool {
	return o.Key != "" && o.Text != ""
}

func (s *SelectionDisplay) IsValid() bool {
	for _, o := range s.Options {
		if !o.IsValid() {
			return false
		}
	}

	return s.Description != ""
}
