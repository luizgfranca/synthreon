package server

import (
	commonmodule "synthreon/modules/common"
	orchestratormodule "synthreon/modules/orchestrator"
	"synthreon/server/api"
)

type AppHandlers struct {
	ProjectAPI          *api.Project
	ToolAPI             *api.Tool
	AuthenticationAPI   *api.Authentication
	OrchestratorService *orchestratormodule.OrchestratorService
	WebHandler          *commonmodule.SPAHandler
}
