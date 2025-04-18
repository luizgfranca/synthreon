package orchestratormodule

import (
	"log"
	clientmodule "synthreon/modules/client"
	configurationmodule "synthreon/modules/configuration"
	projectmodule "synthreon/modules/project"
	providermodule "synthreon/modules/provider"
	sessionmodule "synthreon/modules/session"
	toolmodule "synthreon/modules/tool"
	"synthreon/modules/toolentity"
	tooleventmodule "synthreon/modules/toolevent"
	usermodule "synthreon/modules/user"
)

// TODO: should centralize the concept of context here,
// right now both clientManager and prodiderManager deal with
// htem separately

// FIXME: should find a way to handle when a client connects to a tool
// but there's no provider yet registered for that tool.
//
// My expectation would be that the orchestrator should keep the message
// for some time until the provider connects, if this message expires then it
// drops the event and sinalizes the problem to the frontend
type OrchestratorService struct {
	projectService *projectmodule.ProjectService
	toolService    *toolmodule.ToolService
	configService  *configurationmodule.ConfigurationService

	providerManager *providermodule.ProviderManagerService
	clientManager   *clientmodule.ClientManagerService
}

func NewOrchestratorService(
	projectService *projectmodule.ProjectService,
	toolService *toolmodule.ToolService,
	configService *configurationmodule.ConfigurationService,
) *OrchestratorService {
	o := OrchestratorService{
		projectService: projectService,
		toolService:    toolService,
		configService:  configService,
	}

	o.providerManager = providermodule.NewProviderManagerService(
		&o,
		projectService,
		toolService,
		o.configService.RetryTimeoutSeconds,
        o.configService.AllowToolAutoCreation,
	)

	o.clientManager = clientmodule.NewCLientManagerService(
		&o, projectService,
	)

	return &o
}

func (o *OrchestratorService) RegisterClientEntity(
	session *sessionmodule.Session,
	user *usermodule.User,
	entity toolentity.ToolEntityAdapter,
) {
	if session == nil || user == nil || entity == nil {
		log.Fatalln("[OrchestratorService] null parameters when trying to register client entity", session, user, entity)
	}
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
