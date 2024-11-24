package main

import (
	"log"
	"net/http"
	"platformlab/controlpanel/api"
	"platformlab/controlpanel/api/middleware"
	"platformlab/controlpanel/model"
	"platformlab/controlpanel/service"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	AccessTokenSecretKey string = "supersecret"
)

func DoMigrations(db *gorm.DB) {
	db.AutoMigrate(&model.Project{})
	db.AutoMigrate(&model.Tool{})
	db.AutoMigrate(&model.User{})
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

func CreateDefaultUserIfNotExists(db *gorm.DB) {
	s := service.User{Db: db}

	defaultUser, err := model.NewUser("admin", "test@test.com", "password")
	if err != nil {
		panic(err.Error())
	}

	user, _ := s.FindByEmail(defaultUser.Email)
	if user != nil {
		return
	}

	_, err = s.Create(defaultUser)
	if err != nil {
		panic(err.Error())
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware", r.Method)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
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

	println("asseting creation of default user")
	CreateDefaultUserIfNotExists(db)
	println("done")

	projectAPI := api.ProjectRESTApi(db)
	toolAPI := api.ToolRestAPI(db)
	authenticationAPI := api.AuthenticationRESTApi(db, AccessTokenSecretKey)
	tableAPI := api.Table{}
	sessionMiddleware := middleware.SessionMiddleware(AccessTokenSecretKey)

	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	router.HandleFunc("/api/auth/login", authenticationAPI.Login()).Methods("POST")
	router.HandleFunc("/api/tool/event", toolAPI.GetEventRresponseTEST()).Methods("POST")
	router.HandleFunc("/api/tool/provider/ws", toolAPI.ToolProviderWebsocket()).Methods("GET")
	router.HandleFunc("/api/table", tableAPI.GetTablesMetadata())

	authenticatedRouter := router.PathPrefix("/api").Subrouter()
	authenticatedRouter.Use(sessionMiddleware)
	authenticatedRouter.HandleFunc("/project", projectAPI.GetAllProjects()).Methods("GET")
	authenticatedRouter.HandleFunc("/project", projectAPI.CreateProject()).Methods("POST")
	authenticatedRouter.HandleFunc("/project/{project}/tool", projectAPI.GetToolsFromProject()).Methods("GET")
	authenticatedRouter.HandleFunc("/project/{project}/tool", projectAPI.CreateToolForProject()).Methods("POST")
	authenticatedRouter.HandleFunc("/tool/client/ws", toolAPI.ToolClientWebsocket()).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/control-panel", http.StatusFound)
	})

	router.PathPrefix("/control-panel").Handler(http.StripPrefix("/control-panel", http.FileServer(http.Dir("./web/dist"))))
	router.PathPrefix("/assets").Handler(http.FileServer(http.Dir("./web/dist")))

	println("listening at :8080")
	http.ListenAndServe(":8080", corsMiddleware(router))
}
