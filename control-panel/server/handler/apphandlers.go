package server

import (
	orchestratormodule "platformlab/controlpanel/modules/orchestrator"
	"platformlab/controlpanel/server/api"
)

type AppHandlers struct {
	ProjectAPI          *api.Project
	ToolAPI             *api.Tool
	AuthenticationAPI   *api.Authentication
	OrchestratorService *orchestratormodule.OrchestratorService
}
