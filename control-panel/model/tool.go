package model

type Tool struct {
	ID          uint `json:"id"`
	ProjectId   uint
	Acronym     string `json:"acronym"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
