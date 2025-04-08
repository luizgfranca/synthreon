package route

import (
	"log"
	server "synthreon/server/handler"

	"github.com/gorilla/mux"
)

func SetupBaseRoutes(
	router *mux.Router,
	appHandlers *server.AppHandlers,
) {
	log.Println("[Server][Base] setting up routes...")

	router.Methods("OPTIONS").HandlerFunc(server.CorsOptionsHandler)
	router.HandleFunc("/auth/login", appHandlers.AuthenticationAPI.Login()).Methods("POST")
	router.HandleFunc("/ws/tool/provider", appHandlers.ToolAPI.ToolProviderWebsocket()).Methods("GET")
}
