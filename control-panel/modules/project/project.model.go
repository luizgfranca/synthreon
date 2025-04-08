package projectmodule

import toolmodule "synthreon/modules/tool"

type Project struct {
	ID          uint   `json:"id"`
	Acronym     string `json:"acronym"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Tools []toolmodule.Tool
}

func (p *Project) IsValid() bool {
	return !(p.Acronym == "" || p.Name == "")
}
