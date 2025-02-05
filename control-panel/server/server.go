package server

import (
	"log"
	"net/http"
	commonmodule "platformlab/controlpanel/modules/common"
	configurationmodule "platformlab/controlpanel/modules/configuration"
	"platformlab/controlpanel/server/api"
	"platformlab/controlpanel/server/api/middleware"
	server "platformlab/controlpanel/server/handler"
	"platformlab/controlpanel/server/route"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func StartServer(addr string, configService *configurationmodule.ConfigurationService, db *gorm.DB) {
	appHandlers := &server.AppHandlers{
		ProjectAPI:        api.ProjectRESTApi(db),
		ToolAPI:           api.ToolRestAPI(db),
		AuthenticationAPI: api.AuthenticationRESTApi(db, configService.AccessTokenSecret),
		WebHandler: &commonmodule.SPAHandler{
			StaticPath: configService.StaticFilesDir,
			IndexPath:  "index.html",
		},
	}

	router := mux.NewRouter()
	router.Use(middleware.GetCORSMiddleware())
	route.SetupBaseRoutes(router, appHandlers)
	route.SetupWebRoutes(router, appHandlers)

	authenticatedRouter := router.PathPrefix("/api").Subrouter()
	authenticatedRouter.Use(middleware.GetContentTypeMiddleware(middleware.ContentTypeJSON))
	authenticatedRouter.Use(middleware.GetSessionMiddleware(configService.AccessTokenSecret))
	route.SetupAuthenticatedRoutes(authenticatedRouter, appHandlers)

	log.Println("[Server] listening at", addr)
	http.ListenAndServe(addr, router)
}
