package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"platformlab/controlpanel/model"
	"platformlab/controlpanel/service"
	"platformlab/controlpanel/util"

	"gorm.io/gorm"
)

type Project struct {
	projectService service.Project
}

func (p *Project) GetAllProjects() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		projects := p.projectService.FindAll()
		json.NewEncoder(w).Encode(&projects)
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

func ProjectRESTApi(db *gorm.DB) *Project {
	service := service.Project{Db: db}
	return &Project{projectService: service}
}
