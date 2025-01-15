package providermodule

import (
	commonmodule "platformlab/controlpanel/modules/common"
	projectmodule "platformlab/controlpanel/modules/project"
	toolmodule "platformlab/controlpanel/modules/tool"
	"platformlab/controlpanel/modules/toolentity"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
)

// FIXME: finish implementation

type ProviderManagerService struct {
	// FIXME: reference to orchestrator when tere is one

	projectService *projectmodule.ProjectService

	// FIXME: in the current state, the provider list will allways be growing
	providers []*Provider
}

// DistributeEvent implements Manager.
func (p *ProviderManagerService) DistributeEvent(e *tooleventmodule.ToolEvent) {
	panic("unimplemented: needs orchestrator")
}

// FindProject implements Manager.
func (p *ProviderManagerService) FindProject(acronym string) (*projectmodule.Project, error) {
	maybeProject, err := p.projectService.FindByAcronym(acronym)
	if err != nil {
		return nil, err
	}

	if maybeProject == nil {
		return nil, &commonmodule.GenericLogicError{Message: "project not fonud"}
	}

	return maybeProject, nil
}

// FindTool implements Manager.
func (p *ProviderManagerService) FindTool(project *projectmodule.Project, acronym string) (*toolmodule.Tool, error) {
	maybeTool, err := p.projectService.FindToolByAcronym(project, acronym)
	if err != nil {
		return nil, err
	}
	if maybeTool == nil {
		return nil, &commonmodule.GenericLogicError{Message: "tool from project not fonud"}
	}

	return maybeTool, nil
}

func NewProviderManagerService(
	projectService *projectmodule.ProjectService,
	toolService *toolmodule.ToolService,
) ProviderManagerService {
	return ProviderManagerService{
		projectService: projectService,
		providers:      []*Provider{},
	}
}

func (p *ProviderManagerService) EntityConnection(entity toolentity.ToolEntityAdapter) {
	p.providers = append(p.providers, NewProvider(p, entity))
}
