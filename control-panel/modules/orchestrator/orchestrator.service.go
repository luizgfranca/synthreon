package orchestratormodule

import (
	projectmodule "platformlab/controlpanel/modules/project"
	providermodule "platformlab/controlpanel/modules/provider"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
)

type OrchestratorService struct {
	projectService *projectmodule.ProjectService

	providerManager *providermodule.ProviderManagerService
	// FIXME: client maanger goes here

}

func (o *OrchestratorService) routeEventToProvider(e *tooleventmodule.ToolEvent) {
	panic("uninplemented")
}
