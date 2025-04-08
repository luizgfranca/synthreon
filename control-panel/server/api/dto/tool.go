package api

import (
	toolmodule "synthreon/modules/tool"
)

type CreateToolDto struct {
	Acronym     string `json:"acronym"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *CreateToolDto) IsValid() bool {
	return c.Acronym != "" && c.Name != "" && c.Description != ""
}

func (c *CreateToolDto) ToTool(projectId uint) *toolmodule.Tool {
	return &toolmodule.Tool{
		ProjectId:   projectId,
		Acronym:     c.Acronym,
		Name:        c.Name,
		Description: c.Description,
	}
}
