package clientmodule

import (
	"log"
	commonmodule "platformlab/controlpanel/modules/common"
	projectmodule "platformlab/controlpanel/modules/project"
	sessionmodule "platformlab/controlpanel/modules/session"
	toolmodule "platformlab/controlpanel/modules/tool"
	"platformlab/controlpanel/modules/toolentity"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
	usermodule "platformlab/controlpanel/modules/user"
	"sync"

	"github.com/google/uuid"
)

type Orchestrator interface {
	ForwardEventToProvider(e *tooleventmodule.ToolEvent)
}

// TODO: create abstractions for managerServices
// FIXME: should implement client unregistration
type ClientManagerService struct {
	orchestrator Orchestrator

	contextClientResolver ContextClientResolver

	projectService *projectmodule.ProjectService

	clientsLock sync.Mutex
	clients     []*Client
}

func NewCLientManagerService(
	orchestrator Orchestrator,
	projectService *projectmodule.ProjectService,
) *ClientManagerService {
	return &ClientManagerService{
		orchestrator:          orchestrator,
		projectService:        projectService,
		contextClientResolver: ContextClientResolver{},
		clients:               []*Client{},
	}
}

// FindProject implements Manager.
func (s *ClientManagerService) FindProject(acronym string) (*projectmodule.Project, error) {
	s.log("looking for project: ", acronym)
	maybeProject, err := s.projectService.FindByAcronym(acronym)
	if err != nil {
		return nil, err
	}

	if maybeProject == nil {
		return nil, &commonmodule.GenericLogicError{Message: "project not fonud"}
	}
	project := maybeProject
	s.log("project found: ", project)

	return project, nil
}

// FindTool implements Manager.
func (s *ClientManagerService) FindTool(project *projectmodule.Project, acronym string) (*toolmodule.Tool, error) {
	s.log(
		"looking for referenced tool: \n",
		"project: ", project.Acronym, "\n",
		"tool: ", acronym,
	)

	maybeTool, err := s.projectService.FindToolByAcronym(project, acronym)
	if err != nil {
		return nil, err
	}
	tool := maybeTool

	if tool == nil {
		return nil, &commonmodule.GenericLogicError{Message: "tool from project not fonud"}
	}

	s.log("tool found:", tool)
	return tool, nil
}

// DistributeEvent implements Manager.
func (s *ClientManagerService) DistributeEvent(e *tooleventmodule.ToolEvent) {
	s.orchestrator.ForwardEventToProvider(e)
}

// RegisterClientToolOopen implements Manager.
func (s *ClientManagerService) RegisterClientToolOopen(client *Client) (contextId string) {
	ctxid := uuid.NewString()

	s.log(
		"tool open by client.\n",
		"context: ", ctxid, "\n",
		"client: ", client, "\n",
	)
	ok := s.contextClientResolver.TryRegister(ctxid, client)
	for !ok {
		ok = s.contextClientResolver.TryRegister(ctxid, client)
	}

	return ctxid
}

// TODO: get user from session instead of requiring arg
func (s *ClientManagerService) EntityConnection(
	session *sessionmodule.Session,
	user *usermodule.User,
	entity toolentity.ToolEntityAdapter,
) {
	s.log("new entity connected: ", entity)
	c := NewClient(
		s,
		entity,
		user,
		session,
	)

	s.log("saving new client: ", c.ID)
	s.clientsLock.Lock()
	s.clients = append(s.clients, c)
	s.clientsLock.Unlock()

	c.Start()
}

func (s *ClientManagerService) SendEvent(e *tooleventmodule.ToolEvent) error {
	// BUG: when command/finish without contextId arriving here and crashing the application
	if e.ContextId == "" {
		log.Fatalln("[ClientManagerService] context shound arrive already filled")
	}

	client := s.contextClientResolver.Resolve(e.ContextId)
	if client == nil {
		log.Fatalln("[ClientManagerService] no client found for internal send event")
	}

	s.log("directing event for sending to client", client)
	client.SendEvent(e)

	return nil
}

func (s *ClientManagerService) UnregisterContext(ctxid string) {
	s.log("unregistering context: ", ctxid)
	client := s.contextClientResolver.Resolve(ctxid)
	if client == nil {
		log.Fatalln("[ClientManagerService] unexpected state: trying to unregister unexisting context: ", ctxid)
	}

	client.UnregisterContext(ctxid)

	s.contextClientResolver.Unregister(ctxid)
}

func (s *ClientManagerService) log(v ...any) {
	x := append([]any{"[ClientManagerService]"}, v...)

	log.Println(x...)
}
