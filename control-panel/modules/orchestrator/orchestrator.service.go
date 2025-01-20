package orchestratormodule

import (
	projectmodule "platformlab/controlpanel/modules/project"
	providermodule "platformlab/controlpanel/modules/provider"
	toolmodule "platformlab/controlpanel/modules/tool"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
)

type OrchestratorService struct {
	projectService *projectmodule.ProjectService
	toolService    *toolmodule.ToolService

	providerManager *providermodule.ProviderManagerService
	// FIXME: client maanger goes here

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
	panic("unimplemented")
}

func (o *OrchestratorService) ForwardEventToProvider(e *tooleventmodule.ToolEvent) {
	o.providerManager.SendEvent(e)
}
