package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	api "platformlab/controlpanel/api/dto"
	"platformlab/controlpanel/model"
	"platformlab/controlpanel/service"
	"platformlab/controlpanel/util"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Project struct {
	projectService service.Project
	toolService    service.Tool
}

func (p *Project) getProjectParameter(r *http.Request) (*model.Project, error) {
	params := mux.Vars(r)
	projectAcronym := params["project"]
	if projectAcronym == "" {
		return nil, &util.GenericLogicError{Message: "No project found in request"}
	}

	project, err := p.projectService.FindByAcronym(projectAcronym)
	if err != nil {
		return nil, &util.GenericLogicError{Message: "Project not found"}
	}

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

		foundTools := p.toolService.FindAssociatedWithProject(project)
		json.NewEncoder(w).Encode(&foundTools)
	}
}

func (p *Project) CreateProject() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input model.Project
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			fmt.Println(err.Error())
			json.NewEncoder(w).Encode(ErrorMessage{Message: err.Error()})
			return
		}
		util.Probe(input)

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
	projectService := service.Project{Db: db}
	toolService := service.Tool{Db: db}
	return &Project{projectService: projectService, toolService: toolService}
}
