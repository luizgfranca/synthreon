package toolmodule

type Tool struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	ProjectId   uint   `json:"project_id" gorm:"column:project_id"`
	Acronym     string `json:"acronym"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
