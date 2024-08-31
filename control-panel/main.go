package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Project struct {
	Acronym string `json:"acronym"`
	Name    string `json:"name"`
}

func GetProjects() *[]Project {
	p := []Project{
		{Acronym: "dcc", Name: "DCC"},
		{Acronym: "dsi", Name: "DSI"},
		{Acronym: "customer-identity", Name: "Customer Identity"},
	}

	return &p
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/project", func(w http.ResponseWriter, r *http.Request) {
		var projects = GetProjects()

		json.NewEncoder(w).Encode(projects)
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world\n")
	})

	http.ListenAndServe(":8080", router)
}
