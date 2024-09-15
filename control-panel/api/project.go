package api

import (
	"context"
	"net/http"
	"platformlab/controlpanel/component"
	"platformlab/controlpanel/service"

	"gorm.io/gorm"
)

type Project struct {
	projectService service.Project
}

func (p *Project) GetAllProjects() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		component.ProjectList(*p.projectService.FindAll()).Render(context.Background(), w)
	}
}

func ProjectRESTApi(db *gorm.DB) *Project {
	service := service.Project{Db: db}
	return &Project{projectService: service}
}
