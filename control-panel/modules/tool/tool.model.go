package toolmodule

type Tool struct {
	ID          uint   `json:"id"`
	ProjectId   uint   `json:"project_id"`
	Acronym     string `json:"acronym"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
