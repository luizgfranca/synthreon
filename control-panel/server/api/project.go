package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	commonmodule "synthreon/modules/common"
	projectmodule "synthreon/modules/project"
	toolmodule "synthreon/modules/tool"
	api "synthreon/server/api/dto"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Project struct {
	projectService projectmodule.ProjectService
	toolService    toolmodule.ToolService
}

func (p *Project) getProjectParameter(r *http.Request) (*projectmodule.Project, error) {
	params := mux.Vars(r)
	projectAcronym := params["project"]
	if projectAcronym == "" {
		return nil, &commonmodule.GenericLogicError{Message: "No project found in request"}
	}

	log.Println("looking for project referenced in parameter", projectAcronym)
	project, err := p.projectService.FindByAcronym(projectAcronym)
	if err != nil {
		return nil, &commonmodule.GenericLogicError{Message: "Project not found"}
	}
	log.Println("project referenced in parameters found", project.Acronym, project)

	return project, nil
}

func (p *Project) GetAllProjects() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		projects := p.projectService.FindAll()
		json.NewEncoder(w).Encode(&projects)
	}
}

func (p *Project) GetToolsFromProject() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		project, err := p.getProjectParameter(r)
		if err != nil {
			json.NewEncoder(w).Encode(ErrorMessage{Message: err.Error()})
			return
		}

		foundTools := p.projectService.FindTools(project)
		json.NewEncoder(w).Encode(&foundTools)
	}
}

func (p *Project) CreateProject() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input projectmodule.Project
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			fmt.Println(err.Error())
			json.NewEncoder(w).Encode(ErrorMessage{Message: err.Error()})
			return
		}
		commonmodule.Probe(input)

		if !input.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMessage{Message: "invalid request data"})
			return
		}

		created, err := p.projectService.Create(&input)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMessage{Message: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&created)
	}
}

func (p *Project) CreateToolForProject() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input api.CreateToolDto

		project, err := p.getProjectParameter(r)
		if err != nil {
			json.NewEncoder(w).Encode(ErrorMessage{Message: err.Error()})
			return
		}

		err = json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			json.NewEncoder(w).Encode(ErrorMessage{Message: err.Error()})
			return
		}

		if !input.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMessage{Message: "invalid request data"})
			return
		}

		created, err := p.toolService.Create(input.ToTool(project.ID))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMessage{Message: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&created)
	}
}

func ProjectRESTApi(db *gorm.DB) *Project {
	projectService := projectmodule.ProjectService{Db: db}
	toolService := toolmodule.ToolService{Db: db}
	return &Project{projectService: projectService, toolService: toolService}
}
