package providermodule

import (
	projectmodule "platformlab/controlpanel/modules/project"
	toolmodule "platformlab/controlpanel/modules/tool"
	"platformlab/controlpanel/modules/toolentity"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
)

// FIXME: finish implementation

type ProviderManagerService struct {
	// FIXME: reference to orchestrator when tere is one

	projectService *projectmodule.ProjectService
	toolService    *toolmodule.ToolService

	// FIXME: in the current state, the provider list will allways be growing
	providers []*Provider
}

// DistributeEvent implements Manager.
func (p *ProviderManagerService) DistributeEvent(e *tooleventmodule.ToolEvent) {
	panic("unimplemented")
}

// FindProject implements Manager.
func (p *ProviderManagerService) FindProject(acronym string) (*projectmodule.Project, error) {
	panic("unimplemented")
}

// FindTool implements Manager.
func (p *ProviderManagerService) FindTool(acronym string) (*toolmodule.Tool, error) {
	panic("unimplemented")
}

func NewProviderManagerService(
	projectService *projectmodule.ProjectService,
	toolService *toolmodule.ToolService,
) ProviderManagerService {
	return ProviderManagerService{
		projectService: projectService,
		toolService:    toolService,
		providers:      []*Provider{},
	}
}

func (p *ProviderManagerService) EntityConnection(entity toolentity.ToolEntityAdapter) {
	p.providers = append(p.providers, NewProvider(p, entity))
}
