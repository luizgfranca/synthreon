package server

import "platformlab/controlpanel/server/api"

type AppHandlers struct {
	ProjectAPI        *api.Project
	ToolAPI           *api.Tool
	AuthenticationAPI *api.Authentication
}
