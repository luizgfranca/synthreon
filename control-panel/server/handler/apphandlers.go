package server

import (
	commonmodule "platformlab/controlpanel/modules/common"
	orchestratormodule "platformlab/controlpanel/modules/orchestrator"
	"platformlab/controlpanel/server/api"
)

type AppHandlers struct {
	ProjectAPI          *api.Project
	ToolAPI             *api.Tool
	AuthenticationAPI   *api.Authentication
	OrchestratorService *orchestratormodule.OrchestratorService
	WebHandler          *commonmodule.SPAHandler
}
