package main

import (
	"context"
	"net/http"
	"platformlab/controlpanel/component"
	"platformlab/controlpanel/model"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateMockProjects(db *gorm.DB) {
	db.AutoMigrate(&model.Project{})

	p := []model.Project{
		{Acronym: "dcc", Name: "DCC"},
		{Acronym: "dsi", Name: "DSI"},
		{Acronym: "customer-identity", Name: "Customer Identity"},
	}

	for _, it := range p {
		db.Create(it)
	}
}

func GetAllProjects(db *gorm.DB) *[]model.Project {
	var projects []model.Project

	db.Find(&projects)

	return &projects
}

func main() {
	router := mux.NewRouter()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed connecting to database")
	}

	CreateMockProjects(db)

	router.HandleFunc("/project", func(w http.ResponseWriter, r *http.Request) {
		var projects = GetAllProjects(db)
		component.ProjectList(*projects).Render(context.Background(), w)
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		component := Hello("platformlab")
		component.Render(context.Background(), w)
	})

	http.ListenAndServe(":8080", router)
}
