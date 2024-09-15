package service

import (
	"context"
	"net/http"
	"platformlab/controlpanel/component"
	"platformlab/controlpanel/model"

	"gorm.io/gorm"
)

type Project struct {
	Db *gorm.DB
}

func (p *Project) findAll() *[]model.Project {
	var projects []model.Project

	p.Db.Find(&projects)

	return &projects
}

func (p *Project) GetAllProjects() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		component.ProjectList(*p.findAll()).Render(context.Background(), w)
	}
}
