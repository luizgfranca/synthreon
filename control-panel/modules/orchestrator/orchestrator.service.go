package orchestratormodule

import (
	"log"
	clientmodule "platformlab/controlpanel/modules/client"
	projectmodule "platformlab/controlpanel/modules/project"
	providermodule "platformlab/controlpanel/modules/provider"
	sessionmodule "platformlab/controlpanel/modules/session"
	toolmodule "platformlab/controlpanel/modules/tool"
	"platformlab/controlpanel/modules/toolentity"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
	usermodule "platformlab/controlpanel/modules/user"
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

func (o *OrchestratorService) RegisterClientEntity(
	session *sessionmodule.Session,
	user *usermodule.User,
	entity toolentity.ToolEntityAdapter,
) {
	o.clientManager.EntityConnection(session, user, entity)
}

func (o *OrchestratorService) RegisterProviderEntity(entity toolentity.ToolEntityAdapter) {
	o.providerManager.EntityConnection(entity)
}

// ForwardEventToClient implements providermodule.Orchestrator.
func (o *OrchestratorService) ForwardEventToClient(e *tooleventmodule.ToolEvent) {
	o.log("forwarding event to clientManager: ", e)
	o.clientManager.SendEvent(e)
}

// ForwardEventToProvider implements clientmodule.Orchestrator.
func (o *OrchestratorService) ForwardEventToProvider(e *tooleventmodule.ToolEvent) {
	o.log("forwarding event to providerManager: ", e)
	o.providerManager.SendEvent(e)

	if e.Type == tooleventmodule.EventTypeCommandFinish {
		if e.ContextId == "" {
			log.Fatalln("[OrchestratorService] unexpected state: finish command with no context")
		}

		o.finishContext(e.ContextId)
	}
}

func (o *OrchestratorService) finishContext(ctxid string) {
	o.log("unregistering context: ", ctxid)
	o.clientManager.UnregisterContext(ctxid)
	o.providerManager.UnregisterContext(ctxid)
}

func (o *OrchestratorService) log(v ...any) {
	x := append([]any{"[OrchestratorService]"}, v...)

	log.Println(x...)
}
