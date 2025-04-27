package orchestratormodule

import (
	"log"
	clientmodule "synthreon/modules/client"
	configurationmodule "synthreon/modules/configuration"
	contextmodule "synthreon/modules/context"
	projectmodule "synthreon/modules/project"
	providermodule "synthreon/modules/provider"
	sessionmodule "synthreon/modules/session"
	toolmodule "synthreon/modules/tool"
	"synthreon/modules/toolentity"
	tooleventmodule "synthreon/modules/toolevent"
	tooleventresult "synthreon/modules/toolevent/result"
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

    contextRegistry *contextmodule.ContextRegistry
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

    o.contextRegistry = contextmodule.NewContextRegistry()

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

// NotifyProviderError implements providermodule.Orchestrator, thread-safe
//
// if the referenced contextId does not exist doesn't do anything, but emits
// a warning log. It does not automatically unregisters the context inside the provider, it leaves 
// this to its effective disconnection
func (o *OrchestratorService) NotifyProviderError(ctxId string, error providermodule.ProviderError) {
    o.log("looking for referenced context", ctxId, " to notify error")
    ctx := o.contextRegistry.Get(ctxId)

    // this may happpen when there's a general error and multiple tools in the provider error out
    // FIXME: need to think about this edge case a little bit more
    if ctx == nil {
        o.log("WARNING: referenced context", ctxId, "not found")
    }

    e := o.createFatalErrorEvent(ctx, error)
    o.clientManager.SendEvent(e)
}

// ForwardEventToProvider implements clientmodule.Orchestrator.
func (o *OrchestratorService) ForwardEventToProvider(e *tooleventmodule.ToolEvent) {
	o.log("forwarding event to providerManager: ", e)
	o.providerManager.SendEvent(e)

    if e.Type == tooleventmodule.EventTypeInteractionOpen {
        ctx := contextmodule.Context{
            ID: e.ContextId,
            Project: e.Project, 
            Tool: e.Tool,
        }
        o.contextRegistry.Register(&ctx)
    }

	if e.Type == tooleventmodule.EventTypeCommandFinish {
		if e.ContextId == "" {
			log.Fatalln("[OrchestratorService] unexpected state: finish command with no context")
		}

		o.finishContext(e.ContextId)
	}
}

func (o *OrchestratorService) createFatalErrorEvent(
    ctx *contextmodule.Context, err providermodule.ProviderError,
) *tooleventmodule.ToolEvent {
    o.log("creating fatal error event for context ", ctx, " of type ", err)

    e := tooleventmodule.ToolEvent{
        Type: tooleventmodule.EventTypeCommandFinish,
        Project: ctx.Project,
        Tool: ctx.Tool,
        ContextId: ctx.ID,
        Result: &tooleventresult.ToolEventResult{
            Status: tooleventresult.ToolEventResultStatusFailure,
            Message: "Provider disconnected. Look at the corresponding provider's log.",
        },
    }

    return &e
}

func (o *OrchestratorService) finishContext(ctxid string) {
	o.log("unregistering context: ", ctxid)
	o.clientManager.UnregisterContext(ctxid)
	o.providerManager.UnregisterContext(ctxid)
    o.contextRegistry.Unregister(ctxid)
}

func (o *OrchestratorService) log(v ...any) {
	x := append([]any{"[OrchestratorService]"}, v...)

	log.Println(x...)
}
