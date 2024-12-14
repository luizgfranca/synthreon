package route

import (
	"log"
	server "platformlab/controlpanel/server/handler"

	"github.com/gorilla/mux"
)

func SetupAuthenticatedRoutes(
	router *mux.Router,
	appHandlers *server.AppHandlers,
) {
	log.Println("[Server][Authenticated] setting up routes...")

	router.HandleFunc("/project", appHandlers.ProjectAPI.GetAllProjects()).Methods("GET")
	router.HandleFunc("/project", appHandlers.ProjectAPI.CreateProject()).Methods("POST")
	router.HandleFunc("/project/{project}/tool", appHandlers.ProjectAPI.GetToolsFromProject()).Methods("GET")
	router.HandleFunc("/project/{project}/tool", appHandlers.ProjectAPI.CreateToolForProject()).Methods("POST")
	router.HandleFunc("/tool/client/ws", appHandlers.ToolAPI.ToolClientWebsocket()).Methods("GET")

}
