package orchestratormodule

import (
	clientmodule "platformlab/controlpanel/modules/client"
	projectmodule "platformlab/controlpanel/modules/project"
	providermodule "platformlab/controlpanel/modules/provider"
	toolmodule "platformlab/controlpanel/modules/tool"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
)

// TODO: should centralize the concept of context here,
// right now both clientManager and prodiderManager deal with
// htem separately
type OrchestratorService struct {
	projectService *projectmodule.ProjectService
	toolService    *toolmodule.ToolService

	providerManager *providermodule.ProviderManagerService
	clientManager   *clientmodule.ClientManagerService
}

func NewOrchestratorService(
	projectService *projectmodule.ProjectService,
	toolService *toolmodule.ToolService,
) *OrchestratorService {
	o := OrchestratorService{
		projectService: projectService,
		toolService:    toolService,
	}

	providerManager := providermodule.NewProviderManagerService(
		&o, projectService, toolService,
	)

	o.providerManager = &providerManager
	return &o
}

// ForwardEventToClient implements providermodule.Orchestrator.
func (o *OrchestratorService) ForwardEventToClient(e *tooleventmodule.ToolEvent) {
	o.clientManager.SendEvent(e)
}

// ForwardEventToProvider implements clientmodule.Orchestrator.
func (o *OrchestratorService) ForwardEventToProvider(e *tooleventmodule.ToolEvent) {
	o.providerManager.SendEvent(e)
}
