package main

import (
	"net/http"
	"platformlab/controlpanel/api"
	"platformlab/controlpanel/model"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateMockProjects(db *gorm.DB) {
	db.AutoMigrate(&model.Project{})

	p := []model.Project{
		{Acronym: "proja", Name: "ProjA", Description: "This is the A mock project"},
		{Acronym: "proj-b", Name: "Project B", Description: "This is another example project"},
		{Acronym: "proj-c", Name: "PROJECT C", Description: "This is the C mock project, used to have a basic notion of how this will work and look"},
	}

	for _, it := range p {
		db.Create(&it)
	}
}

func main() {
	router := mux.NewRouter()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed connecting to database")
	}

	println("creating mock projects")
	CreateMockProjects(db)
	println("done")

	projectAPI := api.ProjectRESTApi(db)
	tableAPI := api.Table{}

	router.HandleFunc("/project", projectAPI.GetAllProjects()).Methods("GET")
	router.HandleFunc("/project", projectAPI.CreateProject()).Methods("POST")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/control-panel", http.StatusFound)
	})

	router.PathPrefix("/control-panel").Handler(http.StripPrefix("/control-panel", http.FileServer(http.Dir("./web/dist"))))
	router.PathPrefix("/assets").Handler(http.FileServer(http.Dir("./web/dist")))
	// router.Handle("/assets", http.FileServer(http.Dir("./web/dist/assets")))

	router.HandleFunc("/table", tableAPI.GetTablesMetadata())

	println("listening at :8080")
	http.ListenAndServe(":8080", router)
}
