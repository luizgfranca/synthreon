package providermodule

import (
	"log"
	projectmodule "platformlab/controlpanel/modules/project"
	toolmodule "platformlab/controlpanel/modules/tool"
	"platformlab/controlpanel/modules/toolentity"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
	"sync"

	"github.com/google/uuid"
)

// made this interface instead of a reference to make it easy to test it separately
type Manager interface {
	FindProject(acronym string) (*projectmodule.Project, error)
	FindTool(project *projectmodule.Project, acronym string) (*toolmodule.Tool, error)

	DistributeEvent(e *tooleventmodule.ToolEvent)
}

type ProviderStatus string

const (
	ProviderStatusConnected    ProviderStatus = "CONNECTED"
	ProviderStatusActive       ProviderStatus = "ACTIVE"
	ProviderStatusDisconnected ProviderStatus = "DISCONNECTED"
	ProviderStatusInactive     ProviderStatus = "INACTIVE"
)

type Provider struct {
	ID string

	Project *projectmodule.Project

	handshakeId string

	entity toolentity.ToolEntityAdapter

	handlerWriteLock sync.Mutex
	handlers         []Handler

	// TODO: this is very inneficient, I should think of a better way
	// 		 to avoid half-writes later
	statusLock sync.Mutex
	status     ProviderStatus

	manager Manager
}

// SetStatus @ThreadSafe
func (p *Provider) SetStatus(v ProviderStatus) {
	p.statusLock.Lock()
	defer p.statusLock.Unlock()

	p.status = v
}

// GetStatus @ThreadSafe
func (p *Provider) Status() ProviderStatus {
	p.statusLock.Lock()
	defer p.statusLock.Unlock()

	return p.status
}

func (p *Provider) Disconnect() {
	// FIXME: see #a
	p.SetStatus(ProviderStatusDisconnected)
	p.entity.Close()
}

func (p *Provider) Start() {
	p.log("starting provider")

	p.entity.OnEventReceived(func(e *tooleventmodule.ToolEvent) {
		p.onProviderToServerEvent(e)
	})

	p.entity.OnDisconnect(func() {
		p.SetStatus(ProviderStatusDisconnected)
	})

	p.entity.StartHandler()
}

func NewProvider(
	manager Manager,
	entity toolentity.ToolEntityAdapter,
) *Provider {
	id := uuid.New().String()

	p := Provider{
		ID:       id,
		entity:   entity,
		handlers: []Handler{},
		status:   ProviderStatusConnected,
		manager:  manager,
	}

	p.Start()

	return &p
}

func (p *Provider) log(v ...any) {
	x := append([]any{"[Provider-" + p.ID + "]"}, v...)

	log.Println(x...)
}

func (p *Provider) sendEvent(e *tooleventmodule.ToolEvent) (success bool) {
	p.log("sending event: ", e)
	err := p.entity.SendEvent(e)
	if err != nil {
		p.log("error sending event", err.Error())
		p.Disconnect()
		return false
	}

	return true
}

func (p *Provider) completeHandshake(project *projectmodule.Project) {
	handshakeId := uuid.NewString()

	// i only lock the status because any operation on the other items will look for it first
	// the only important problem here is some thread reading the status and assumimng everything
	// is ready because it is active, but the looking for the project and not finding it
	p.statusLock.Lock()
	defer p.statusLock.Unlock()

	p.Project = project
	p.handshakeId = handshakeId
	p.status = ProviderStatusActive
}

func (p *Provider) createHandler(t *toolmodule.Tool) *Handler {
	handler := NewHandler(t)

	p.handlerWriteLock.Lock()
	defer p.handlerWriteLock.Unlock()

	p.handlers = append(p.handlers, handler)

	return &handler
}

// FIXME: i should handle when the provider tries to reconnect here
func (p *Provider) onProviderToServerEvent(e *tooleventmodule.ToolEvent) {
	switch p.Status() {
	case ProviderStatusConnected:
		if !e.IsHandshake() {
			// FIXME: make a message to indicate to the client what happened
			// 		  and maybe not disconnect them immediately
			p.log("error: invalid first message from new connection: ", e.Type)
			p.Disconnect()
			return
		}
		p.handleHandshakeEvent(e)
	case ProviderStatusActive:
		if e.IsHandshake() {
			// FIXME: should make a way to handle this invalid stage problems
			p.log("error: handshake stage request received while active: ", e.Type)
			p.Disconnect()
			return
		}

		if e.IsAnnouncement() {
			p.handleAnnouncementEvent(e)
			return
		}

		p.handleActiveProviderEvent(e)
	case ProviderStatusDisconnected:
		// FIXME: I may not panic here, if we are between marking the status as
		// 		  disconnected and effectively disconnecting from the entity this
		// 		  case may occur.
		// 		  The basic idea would be to just drop the event, but there may be
		// 		  a better solution
		// 		  #a
		log.Fatal("[Provider] unexpected state: event received from provider that should be disconnected: ", e)
	}
}

func (p *Provider) handleHandshakeEvent(e *tooleventmodule.ToolEvent) {
	if p.Status() != ProviderStatusConnected {
		panic("unexpected provider status for handleHandshakeEvent")
	}

	if !e.IsHandshake() {
		panic("unexpected event status for handleHandshakeEvent")
	}

	if e.Type != tooleventmodule.EventTypeHandshakeRequest {
		// FIXME: should make a way to handle this invalid stage problems
		p.log("error: event is not provider->server: ", e.Type)
		p.Disconnect()
		return
	}

	project, err := p.manager.FindProject(e.Project)
	if err != nil {
		p.log("project on handshake request not found: ", e.Project)
		nack := tooleventmodule.ToolEvent{
			Type:    tooleventmodule.EventTypeHandshakeNACK,
			Project: e.Project,
			Reason:  "project.invalid",
		}

		if ok := p.sendEvent(&nack); !ok {
			p.Disconnect()
			return
		}

		return
	}

	p.completeHandshake(project)
	ack := tooleventmodule.ToolEvent{
		Type:        tooleventmodule.EventTypeHandshakeACK,
		Project:     p.Project.Acronym,
		HandshakeId: p.handshakeId,
		ProviderId:  p.ID,
	}
	p.sendEvent(&ack)
}

func (p *Provider) handleAnnouncementEvent(e *tooleventmodule.ToolEvent) {
	if p.Status() != ProviderStatusActive {
		panic("unexpected provider status for handleHandshakeEvent")
	}

	if !e.IsAnnouncement() {
		panic("unexpected event status for handleHandshakeEvent")
	}

	if e.Type != tooleventmodule.EventTypeAnnouncementHandler {
		// FIXME: should make a way to handle this invalid stage problems
		p.log("error: event is not provider->server: ", e.Type)
		p.Disconnect()
		return
	}

	if e.Project != p.Project.Acronym {
		log.Printf(`[Provider] error: project from request differs from provider's registered: 
			(provider's: %s, request's: %s)`,
			p.Project.Acronym, e.Project,
		)

		nack := tooleventmodule.ToolEvent{
			Type:        tooleventmodule.EventTypeAnnouncementNACK,
			Project:     p.Project.Acronym,
			Tool:        e.Tool,
			HandshakeId: p.handshakeId,
			ProviderId:  p.ID,
			Reason:      "project.match.none",
		}
		if ok := p.sendEvent(&nack); !ok {
			p.Disconnect()
			return
		}

		return
	}

	tool, err := p.manager.FindTool(p.Project, e.Tool)
	if err != nil {
		p.log("error: provider tried to announce invalid tool", e.Tool)
		nack := tooleventmodule.ToolEvent{
			Type:        tooleventmodule.EventTypeAnnouncementNACK,
			Project:     p.Project.Acronym,
			Tool:        e.Tool,
			HandshakeId: p.handshakeId,
			ProviderId:  p.ID,
			Reason:      "tool.invalid",
		}
		if ok := p.sendEvent(&nack); !ok {
			p.log("not ok")
			p.Disconnect()
			return
		}

		return
	}

	created := p.createHandler(tool)
	ack := tooleventmodule.ToolEvent{
		Type:        tooleventmodule.EventTypeAnnouncementACK,
		Project:     p.Project.Acronym,
		Tool:        tool.Acronym,
		HandshakeId: p.handshakeId,
		ProviderId:  p.ID,
		HandlerId:   created.ID,
	}

	p.sendEvent(&ack)
}

// TODO: there are not tests to stress these validations yet, i should create it later
func (p *Provider) isEventValidForDistribution(e *tooleventmodule.ToolEvent) bool {
	if !e.IsCommand() {
		return false
	}

	if !(e.Project == p.Project.Acronym &&
		e.ProviderId == p.ID &&
		e.HandlerId != "") {
		return false
	}

	for _, it := range p.handlers {
		if it.ID == e.HandlerId {
			return e.Tool == it.Tool.Acronym
		}
	}

	return false
}

func (p *Provider) handleActiveProviderEvent(e *tooleventmodule.ToolEvent) {
	if !p.isEventValidForDistribution(e) {
		// TODO: just drops for now, maybe it should have a better behavior
		return
	}

	p.manager.DistributeEvent(e)
}
