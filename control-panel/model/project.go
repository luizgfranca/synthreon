package model

type Project struct {
	ID          uint   `json:"id"`
	Acronym     string `json:"acronym"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (p *Project) IsValid() bool {
	return !(p.Acronym == "" || p.Name == "")
}
