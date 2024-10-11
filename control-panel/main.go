package main

import (
	"log"
	"net/http"
	"platformlab/controlpanel/api"
	"platformlab/controlpanel/model"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DoMigrations(db *gorm.DB) {
	db.AutoMigrate(&model.Project{})
	db.AutoMigrate(&model.Tool{})
}

func CreateMockProjects(db *gorm.DB) {
	p := []model.Project{
		{Acronym: "proja", Name: "ProjA", Description: "This is the A mock project"},
		{Acronym: "proj-b", Name: "Project B", Description: "This is another example project"},
		{Acronym: "proj-c", Name: "PROJECT C", Description: "This is the C mock project, used to have a basic notion of how this will work and look"},
	}

	for _, it := range p {
		db.Create(&it)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware", r.Method)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers:", "Origin, Content-Type, X-Auth-Token, Authorization")
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed connecting to database")
	}

	println("doing database migrations")
	DoMigrations(db)
	println("done")

	println("creating mock projects")
	CreateMockProjects(db)
	println("done")

	projectAPI := api.ProjectRESTApi(db)
	tableAPI := api.Table{}

	router.HandleFunc("/api/project", projectAPI.GetAllProjects()).Methods("GET")
	router.HandleFunc("/api/project", projectAPI.CreateProject()).Methods("POST")
	router.HandleFunc("/api/project/{project}/tool", projectAPI.GetToolsFromProject()).Methods("GET")
	router.HandleFunc("/api/project/{project}/tool", projectAPI.CreateToolForProject()).Methods("POST")
	router.HandleFunc("/api/table", tableAPI.GetTablesMetadata())

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/control-panel", http.StatusFound)
	})

	router.PathPrefix("/control-panel").Handler(http.StripPrefix("/control-panel", http.FileServer(http.Dir("./web/dist"))))
	router.PathPrefix("/assets").Handler(http.FileServer(http.Dir("./web/dist")))

	println("listening at :8080")
	http.ListenAndServe(":8080", corsMiddleware(router))
}
