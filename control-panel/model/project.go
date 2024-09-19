package model

type Project struct {
	Acronym string `json:"acronym"`
	Name    string `json:"name"`
}

func (p *Project) IsValid() bool {
	return !(p.Acronym == "" || p.Name == "")
}
