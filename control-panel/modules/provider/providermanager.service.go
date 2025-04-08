package providermodule

import (
	"log"
	commonmodule "synthreon/modules/common"
	projectmodule "synthreon/modules/project"
	toolmodule "synthreon/modules/tool"
	"synthreon/modules/toolentity"
	tooleventmodule "synthreon/modules/toolevent"
	"time"
)

type Orchestrator interface {
	ForwardEventToClient(e *tooleventmodule.ToolEvent)
}

// FIXME: add project and tool deregistration handling
// TODO: create abstractions for managerServices

// FIXME: when a provider disconnects it should notify the orchestrator in order for
// it to cancel the context
type ProviderManagerService struct {
	orchestrator Orchestrator

	projectService *projectmodule.ProjectService
	toolService    *toolmodule.ToolService

	contextProviderResolver        ContextProviderResolver
	projectAndToolProviderResolver ProjectAndToolProviderResolver

	retryTimeoutSeconds  int
	shouldAutoCreateTool bool
	dlq                  *EventDLQ

	// TODO: is this list irrelevant?
	providers []*Provider
}

// OnProviderDisconnection implements Manager.
func (p *ProviderManagerService) OnProviderDisconnection(provider *Provider) {
	if provider == nil {
		log.Fatalln("provider to be unregistered from providermanager should not be null")
	}

	p.log("provider ", provider.ID, " disconnected, degeristering mappings")
	p.projectAndToolProviderResolver.UnregisterProviderEntries(provider)
	p.contextProviderResolver.UnregisterProviderEntries(provider)

	// TODO: could this have a race condition if the provider unregisters right after connecting?
	for i, it := range p.providers {
		if it.ID == provider.ID {
			commonmodule.RemoveFromUnorderedSlice(p.providers, i)
			break // provider should not have duplicates in the list
		}
	}
}

// NewProviderManagerService:
// orchestrator, projectService and toolServie are not nullable
// retryTimeoutSeconds default = 0 (no retry)
func NewProviderManagerService(
	orchestrator Orchestrator,
	projectService *projectmodule.ProjectService,
	toolService *toolmodule.ToolService,
	retryTimeoutSeconds int,
	shouldAutoCreateTool bool,
) *ProviderManagerService {
	if orchestrator == nil || projectService == nil || toolService == nil {
		log.Fatalln("tried to start ProviderManager with null dependency", orchestrator, projectService, toolService)
	}

	p := ProviderManagerService{
		orchestrator:         orchestrator,
		projectService:       projectService,
		toolService:          toolService,
		retryTimeoutSeconds:  retryTimeoutSeconds,
		shouldAutoCreateTool: shouldAutoCreateTool,

		dlq:                            newEventDLQ(retryTimeoutSeconds),
		contextProviderResolver:        ContextProviderResolver{},
		projectAndToolProviderResolver: ProjectAndToolProviderResolver{},
		providers:                      []*Provider{},
	}

	// TODO: Should i really use a daemon here? This was simpler when only the providers ran on the
	// background
	go p.dlq.sweepService()

	return &p
}

// DistributeEvent implements Manager.
// called by the managed providers
func (p *ProviderManagerService) DistributeEvent(e *tooleventmodule.ToolEvent) {
	if e.Type == tooleventmodule.EventTypeCommandFinish {
		p.contextProviderResolver.Unregister(e.ContextId)
	}

	e.HandshakeId = ""
	e.ExecutionId = ""

	p.orchestrator.ForwardEventToClient(e)
}

// retryFailedEventsFromProjectAndTool:
// presumes already verified project and tool acronyms,
// failed retries are readded to queue by underlying behavior
func (p *ProviderManagerService) retryFailedEventsFromProjectAndTool(projectAcronym string, toolAcronym string) {
	// FIXME: needs to wait for a bit because the internal registration happens
	// before the ACK message is sent to the provider
	// instead of this, which is error prone, there are two options to think about:
	// - a status for the handlers that should be controlled here for when their registration if fully done
	// - consider registering the handler only after the ACK message is sent
	time.Sleep(1 * time.Second)
	p.log("retrying failed events")

	eventToRetry := p.dlq.popFromProjectAndTool(projectAcronym, toolAcronym)
	for eventToRetry != nil {
		p.log("retrying event", eventToRetry)
		// we don't save the event again in teh DLQ here because there's already a logic to
		// put it there again in the send event logic,
		// this function assumes that it will behave the same as if it was a new event send.
		p.SendEvent(eventToRetry)

		eventToRetry = p.dlq.popFromProjectAndTool(projectAcronym, toolAcronym)
	}
}

// RegisterProviderProjectAndTool implements Manager.
// m is not nullable,
func (p *ProviderManagerService) RegisterProviderProjectAndTool(m *ProviderToolMapping) {
	if m == nil {
		log.Fatalln("in manager provider registration, but provider mapping is null")
	}

	p.projectAndToolProviderResolver.Register(m.Project.Acronym, m.Tool.Acronym, m.Provider)
	go p.retryFailedEventsFromProjectAndTool(m.Project.Acronym, m.Tool.Acronym)
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
	p.log(
		"looking for announced tool: \n",
		"project: ", project.Acronym, "\n",
		"tool: ", acronym,
	)

	maybeTool, err := p.projectService.FindToolByAcronym(project, acronym)
	if err != nil {
		return nil, err
	}

	if maybeTool == nil {
		return nil, &commonmodule.GenericLogicError{Message: "tool from project not fonud"}
	}

	p.log("tool found:", maybeTool)
	return maybeTool, nil
}

func (p *ProviderManagerService) TryCreateTool(
	project *projectmodule.Project,
	acronym string,
	props *tooleventmodule.ToolProperties,
) (*toolmodule.Tool, error) {
	if !p.shouldAutoCreateTool {
		p.log("ALLOW_TOOL_AUTOCREATION flag unmarked, will not create new tool")
		return nil, &commonmodule.GenericLogicError{Message: "should not autocreate non existing tools"}
	}

	t := toolmodule.Tool{
		ProjectId:   project.ID,
		Acronym:     acronym,
		Name:        acronym,
		Description: "",
	}

	if props != nil {
		t.Name = props.Name
		t.Description = props.Description
	}

	p.log("trying to create tool: ", t)
	tool, err := p.toolService.Create(&t)
	if err != nil {
		return nil, err
	}

	return tool, nil
}

func (p *ProviderManagerService) EntityConnection(entity toolentity.ToolEntityAdapter) {
	p.providers = append(p.providers, NewProvider(p, entity))
}

func (p *ProviderManagerService) SendEvent(e *tooleventmodule.ToolEvent) error {
	p.log("event received by providerManager")
	if e.ContextId == "" || e.Tool == "" || e.Project == "" {
		log.Fatalln("[ProviderManagerService] unexptected event attributes when reaching providerManager", e)
	}

	p.log("trying to route")
	err := p.contextProviderResolver.TryRouteEvent(e)
	if err != nil {
		switch err.(type) {
		case *ContextNotFounError:
			p.log("context not found, creating")
			provider, err := p.projectAndToolProviderResolver.Resolve(e.Project, e.Tool)
			if err != nil {
				p.log("could not resolve by project and tool, adding event to dead letter queue")

				// this function is also called on retries, this should be considered when doning this logic
				// in this case, this behavior already re-registers the event in the dead letter queue
				// if the retry fails
				// in this registration, the event is saved with a new expiration, i don't think this is a
				// problem, but it should also be considered
				// if this is removed, "retryFailedEventsFromProjectAndTool" or its current equivalent should be
				// reviewed
				p.dlq.register(e)

				// Not sure about this behavior,
				// In this case for example, i don't think it should be an error yet because
				// there's a chance that the event will be reprocessed and will be fine.
				// But should the orchestrator have the knowlege of this??
				// TODO: should reevaluate when thinking of a way to handle geenral errros
				// in event communication.
				return nil
			}

			p.log("registering new context")
			p.contextProviderResolver.Register(e.ContextId, provider)

			p.log("routing event to provider")
			err = p.contextProviderResolver.TryRouteEvent(e)
			if err != nil {
				log.Fatalln("unexpected behavior: could not route event even after context setup")
			}

			return nil
		case *commonmodule.GenericLogicError:
			return err
		}
	}

	return nil
}

func (s *ProviderManagerService) UnregisterContext(ctxid string) {
	s.log("unregistering context: ", ctxid)
	provider := s.contextProviderResolver.TryResolve(ctxid)
	if provider == nil {
		log.Fatalln("[ProviderManagerService] unexpected: trying to unregister unexisting context: ", ctxid)
	}

	provider.UnregisterContext(ctxid)
	s.contextProviderResolver.Unregister(ctxid)
}

func (p *ProviderManagerService) log(v ...any) {
	x := append([]any{"[ProviderManagerService]"}, v...)

	log.Println(x...)
}

// TODO: I dont know if i like creating this struct just for this,
// maybe there's a better approach
type ProviderToolMapping struct {
	Provider *Provider
	Project  *projectmodule.Project
	Tool     *toolmodule.Tool
}
