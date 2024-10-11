package api

import "platformlab/controlpanel/model"

type CreateToolDto struct {
	Acronym     string `json:"acronym"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *CreateToolDto) IsValid() bool {
	return c.Acronym != "" && c.Name != "" && c.Description != ""
}

func (c *CreateToolDto) ToTool(projectId uint) *model.Tool {
	return &model.Tool{
		ProjectId:   projectId,
		Acronym:     c.Acronym,
		Name:        c.Name,
		Description: c.Description,
	}
}
