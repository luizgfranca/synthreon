package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Project struct {
	Acronym string `json:"acronym"`
	Name    string `json:"name"`
}

func CreateMockProjects(db *gorm.DB) {
	db.AutoMigrate(&Project{})

	p := []Project{
		{Acronym: "dcc", Name: "DCC"},
		{Acronym: "dsi", Name: "DSI"},
		{Acronym: "customer-identity", Name: "Customer Identity"},
	}

	for _, it := range p {
		db.Create(it)
	}
}

func GetAllProjects(db *gorm.DB) *[]Project {
	var projects []Project

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

		json.NewEncoder(w).Encode(projects)
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world\n")
	})

	http.ListenAndServe(":8080", router)
}
