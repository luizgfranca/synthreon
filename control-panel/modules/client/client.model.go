package clientmodule

import (
	"log"
	commonmodule "platformlab/controlpanel/modules/common"
	projectmodule "platformlab/controlpanel/modules/project"
	toolmodule "platformlab/controlpanel/modules/tool"
	"platformlab/controlpanel/modules/toolentity"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
	usermodule "platformlab/controlpanel/modules/user"
	"sync"

	"github.com/google/uuid"
)

type ClientStatus string

const (
	ClientStatusActive   ClientStatus = "ACTIVE"
	ClientStatusInactive ClientStatus = "INACTIVE"
)

type Manager interface {
	FindProject(acronym string) (*projectmodule.Project, error)
	FindTool(project *projectmodule.Project, acronym string) (*toolmodule.Tool, error)

	DistributeEvent(e *tooleventmodule.ToolEvent)
}

type Client struct {
	ID                      string
	manager                 Manager
	user                    *usermodule.User
	entity                  toolentity.ToolEntityAdapter
	sessionId               string
	terminals               sync.Map
	contextTerminalResolver ContextTerminalResolver
}

func NewClient(
	manager Manager,
	entity toolentity.ToolEntityAdapter,
	user *usermodule.User,
	sessionId string,
) *Client {
	id := uuid.New().String()

	c := Client{
		ID:                      id,
		manager:                 manager,
		entity:                  entity,
		user:                    user,
		sessionId:               sessionId,
		terminals:               sync.Map{},
		contextTerminalResolver: ContextTerminalResolver{},
	}

	return &c
}

func (c *Client) Start() {
	c.log("starting client")

	c.entity.OnEventReceived(func(e *tooleventmodule.ToolEvent) {
		c.onClientEvent(e)
	})

	c.entity.StartHandler()
}

func (c *Client) SendEvent(e *tooleventmodule.ToolEvent) {
	c.log("sending event: ", e)

	if e.ContextId == "" {
		log.Fatalln("[Client] context empty when sending message to client")
	}

	maybeTerm := c.contextTerminalResolver.Resolve(e.ContextId)
	if maybeTerm == nil {
		c.log("error: (INTERNAL DROP) no terminal found fot context: ", e.ContextId)
		return
	}
	term := maybeTerm

	e.TerminalId = term.ID
	e.SessionId = c.sessionId

	c.entity.SendEvent(e)
}

func (c *Client) openTerminal(projectAcronym string, toolAcronym string) (*Terminal, error) {
	project, err := c.manager.FindProject(projectAcronym)
	if err != nil {
		// FIXME: should have a way to sinalize to client that open errored
		return nil, &commonmodule.GenericLogicError{Message: "could not find referenced project: " + projectAcronym}
	}

	tool, err := c.manager.FindTool(project, toolAcronym)
	if err != nil {
		// FIXME: should have a way to sinalize to client that open errored
		return nil, &commonmodule.GenericLogicError{Message: "could not find referenced tool: " + toolAcronym}
	}

	terminal := Terminal{
		ID:      uuid.NewString(),
		Client:  c,
		Project: project,
		Tool:    tool,
		Status:  TerminalStatusRunning,
	}

	return &terminal, nil
}

func (c *Client) saveTerminalIfDoesntExist(t *Terminal) (success bool) {
	return c.terminals.CompareAndSwap(t.ID, nil, t)
}

func (c *Client) getTerminal(id string) *Terminal {
	v, ok := c.terminals.Load(id)
	if !ok {
		return nil
	}

	t, ok := v.(*Terminal)
	if !ok {
		panic("only pointers to termnals should have been in terminals list")
	}

	return t
}

// TODO: should think on the behavior when the client already informas a terminal
func (c *Client) onOpenEvent(e *tooleventmodule.ToolEvent) {
	c.log("creating new terminal")
	terminal, err := c.openTerminal(e.Project, e.Tool)
	if err != nil {
		// FIXME: should create a way to communicate this to the client
		c.log("error openning terminal: ", err)
		return
	}

	c.log("registering new terminal: ", terminal.ID, terminal)
	success := c.saveTerminalIfDoesntExist(terminal)
	if !success {
		panic("unexpected state, tried to create terminal but it already exists")
	}

	e.SessionId = ""

	c.manager.DistributeEvent(e)
}

func (c *Client) onSImpleEvent(e *tooleventmodule.ToolEvent) {
	if e.TerminalId == "" {
		// FIXME: should create a way to pass this information to the client
		c.log("terminalId is required")
		return
	}

	c.log("looking for terminal")
	term := c.getTerminal(e.TerminalId)
	if term == nil {
		// FIXME: should create a way to pass this information to the client
		c.log("terminal not found: ", e.TerminalId)
		return
	}

	if e.Project != term.Project.Acronym || e.Tool != term.Tool.Acronym {
		// FIXME: should create a way to pass this information to the client
		c.log(
			"terminal and event information do not match: \n",
			"event: ", e, "\n",
			"terminal: ", term, "\n",
		)
		return
	}

	if term.Status != TerminalStatusRunning {
		// FIXME: should create a way to pass this information to the client
		c.log("terminal on invalid status: ", term.Status)
		return
	}

	e.SessionId = ""

	c.manager.DistributeEvent(e)
}

func (c *Client) onClientEvent(e *tooleventmodule.ToolEvent) {
	c.log("event received", e)

	switch e.Type {
	case tooleventmodule.EventTypeInteractionOpen:
		c.onOpenEvent(e)
	default:
		c.onSImpleEvent(e)
	}
}

func (c *Client) log(v ...any) {
	x := append([]any{"[Client-" + c.ID + "]"}, v...)

	log.Println(x...)
}
