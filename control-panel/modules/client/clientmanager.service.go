package clientmodule

import (
	projectmodule "platformlab/controlpanel/modules/project"
	toolmodule "platformlab/controlpanel/modules/tool"
	"platformlab/controlpanel/modules/toolentity"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
)

type Orchestrator interface {
	ForwardEventToClient(e *tooleventmodule.ToolEvent)
}

type ClientManagerService struct {
	orchestrator Orchestrator

	contextClientResolver ContextClientResolver

	projectService *projectmodule.ProjectService

	clients []Client
}

func NewCLientManagerService(
	orchestrator Orchestrator,
	projectService *projectmodule.ProjectService,
) ClientManagerService {
	return ClientManagerService{
		orchestrator:          orchestrator,
		projectService:        projectService,
		contextClientResolver: ContextClientResolver{},
		clients:               []Client{},
	}
}

// FindProject implements Manager.
func (s *ClientManagerService) FindProject(acronym string) (*projectmodule.Project, error) {
	panic("uninplemented")
}

// FindTool implements Manager.
func (s *ClientManagerService) FindTool(project *projectmodule.Project, acronym string) (*toolmodule.Tool, error) {
	panic("uninplemented")
}

// DistributeEvent implements Manager.
func (s *ClientManagerService) DistributeEvent(e *tooleventmodule.ToolEvent) {
	panic("uninplemented")
}

// RegisterClientContext implements Manager.
func (s *ClientManagerService) RegisterContextClient(contextId string, client *Client) {
	panic("uninplemented")
}

func (s *ClientManagerService) EntityConnection(entity toolentity.ToolEntityAdapter) {
	panic("unimplemented")
}

func (s *ClientManagerService) SendEvent(e *tooleventmodule.ToolEvent) error {
	panic("unimplemented")
}
