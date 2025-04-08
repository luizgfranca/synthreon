package route

import (
	"log"
	"net/http"
	server "synthreon/server/handler"

	"github.com/gorilla/mux"
)

func SetupWebRoutes(
	router *mux.Router,
	appHandlers *server.AppHandlers,
) {
	log.Println("[Server][Base] setting up routes...")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/control-panel", http.StatusFound)
	})

	router.PathPrefix("/control-panel").Handler(
		http.StripPrefix(
			"/control-panel",
			appHandlers.WebHandler,
		),
	)
	router.PathPrefix("/assets").Handler(http.FileServer(http.Dir("./web/dist")))
}
