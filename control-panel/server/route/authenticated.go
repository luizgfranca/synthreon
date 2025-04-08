package route

import (
	"log"
	server "synthreon/server/handler"

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

	// FIXME: should generate a temporary one-time use token for here, or authenticate inside the websocket protocol
	// instead of passing the session
	router.HandleFunc("/tool/client/ws/{accessToken}", appHandlers.ToolAPI.ToolClientWebsocket()).Methods("GET")
}
