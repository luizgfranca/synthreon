package providermodule

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	commonmodule "platformlab/controlpanel/modules/common"
	projectmodule "platformlab/controlpanel/modules/project"
	toolmodule "platformlab/controlpanel/modules/tool"
	"platformlab/controlpanel/modules/toolentity"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// TODO: in the future, to be more formal with the tests,
// 		 I may also mock the entity, but i can't bother right
// 		 now :D

type negotionationTest struct {
	Event   tooleventmodule.ToolEvent
	Success bool
	ACK     bool
}

// TODO: make tests with
// - multiple handlers
// - event concurrency
// - events in for incorrect states
func TestProviderRegistrationRules(t *testing.T) {
	exampleValidHandshake := tooleventmodule.ToolEvent{
		Type:    tooleventmodule.EventTypeHandshakeRequest,
		Project: "validproject",
	}

	exampleEventToDistribute := tooleventmodule.ToolEvent{
		Type:        tooleventmodule.EventTypeCommandFinish,
		Project:     "validproject",
		Tool:        "validtool",
		ProviderId:  "", // to be filled
		HandlerId:   "", // to be filled
		ContextId:   uuid.NewString(),
		ExecutionId: uuid.NewString(),
	}

	handshakeBaseTests := []negotionationTest{
		{
			Event: tooleventmodule.ToolEvent{
				Type:    tooleventmodule.EventTypeAnnouncementHandler,
				Project: "validproject",
			},
			Success: false,
			ACK:     false,
		},
		{
			Event:   exampleValidHandshake,
			Success: true,
			ACK:     true,
		},
		{
			Event: tooleventmodule.ToolEvent{
				Type:    tooleventmodule.EventTypeHandshakeRequest,
				Project: "invalidproject",
			},
			Success: true,
			ACK:     false,
		},
	}

	announcementBaseTests := []negotionationTest{
		{
			Event: tooleventmodule.ToolEvent{
				Type:        tooleventmodule.EventTypeAnnouncementACK,
				Project:     "validproject",
				Tool:        "validtool",
				HandshakeId: "", // to be fiiled
				ProviderId:  "", // to be fiiled
			},
			Success: false,
			ACK:     false,
		},
		{
			Event: tooleventmodule.ToolEvent{
				Type:        tooleventmodule.EventTypeAnnouncementHandler,
				Project:     "validproject",
				Tool:        "validtool",
				HandshakeId: "", // to be fiiled
				ProviderId:  "", // to be fiiled
			},
			Success: true,
			ACK:     true,
		},
		{
			Event: tooleventmodule.ToolEvent{
				Type:        tooleventmodule.EventTypeAnnouncementHandler,
				Project:     "validproject",
				Tool:        "invalidtool",
				HandshakeId: "", // to be fiiled
				ProviderId:  "", // to be fiiled
			},
			Success: true,
			ACK:     false,
		},
		{
			Event: tooleventmodule.ToolEvent{
				Type:        tooleventmodule.EventTypeAnnouncementHandler,
				Project:     "validproject",
				Tool:        "invalidtool",
				HandshakeId: "", // to be fiiled
				ProviderId:  "", // to be fiiled
			},
			Success: true,
			ACK:     false,
		},
	}

	manager := testManager{}
	var provider *Provider

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		websocketUpgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
			return true
		}}
		conn, err := websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal("internal test error: ", err.Error())
		}

		entity := toolentity.NewWebsocketToolEntity(conn)
		if entity == nil {
			t.Error("null websocket entity after creation")
		}

		provider = NewProvider(&manager, entity)
	}))
	defer s.Close()

	for i := range handshakeBaseTests {
		it := handshakeBaseTests[i]
		name := fmt.Sprintln("handshake basic test: ", it)
		t.Run(name, func(t *testing.T) {
			fmt.Print(name)
			websocketURL := "ws" + strings.TrimPrefix(s.URL, "http")
			fmt.Println("websocket URL: ", websocketURL)
			ws, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
			if err != nil {
				t.Fatal("internal test error: unable to connect to test websocket: ", err.Error())
			}

			msg, err := tooleventmodule.WriteV0EventString(&it.Event)
			if err != nil {
				t.Fatal("internal test error: unable to encode event: ", it.Event)
			}

			ws.WriteMessage(websocket.TextMessage, []byte(*msg))
			ws.SetReadDeadline(time.Now().Add(time.Second))
			_, buff, err := ws.ReadMessage()
			if err != nil {
				fmt.Println("read error: ", err.Error())
				if !it.Success {
					return
				}

				t.Fatal("internal test error reading websocket message", err.Error())
				return
			}

			str := string(buff)
			event, err := tooleventmodule.ParseEventString(&str)
			if err != nil {
				t.Error("malformed response event", err.Error())
				return
			}

			log.Println("ack validation")
			if it.ACK && event.Type != tooleventmodule.EventTypeHandshakeACK {
				t.Error("expected ACK, got: ", event)
				return
			}

			log.Println("nack validation")
			if !it.ACK && event.Type != tooleventmodule.EventTypeHandshakeNACK {
				t.Error("expected NACK, got: ", event)
				return
			}

			log.Println("should be ack")
			if !it.ACK {
				return
			}

			log.Println("looking for project")
			if provider.Project.Acronym != it.Event.Project || provider.Project.Name == "" {
				t.Error("project not saved correctly on provider: ", provider.Project)
				return
			}
		})
	}

	for i := range announcementBaseTests {
		it := announcementBaseTests[i]
		name := fmt.Sprintln("handshake basic test: ", it)
		t.Run(name, func(t *testing.T) {
			fmt.Print(name)
			websocketURL := "ws" + strings.TrimPrefix(s.URL, "http")
			fmt.Println("websocket URL: ", websocketURL)
			ws, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
			if err != nil {
				t.Fatal("internal test error: unable to connect to test websocket: ", err.Error())
			}

			msg, err := tooleventmodule.WriteV0EventString(&exampleValidHandshake)
			if err != nil {
				t.Fatal("internal test error: unable to encode event: ", it.Event)
			}

			ws.WriteMessage(websocket.TextMessage, []byte(*msg))
			ws.SetReadDeadline(time.Now().Add(time.Second))
			_, buff, err := ws.ReadMessage()
			if err != nil {
				t.Fatal("intenral test error: read error on handshake: ", err.Error())
			}

			str := string(buff)
			ack, err := tooleventmodule.ParseEventString(&str)
			if err != nil {
				t.Error("internal test error: malformed response event", err.Error())
				return
			}

			if it.ACK && ack.Type != tooleventmodule.EventTypeHandshakeACK {
				t.Error("internal test error: handshake not accepted")
				return
			}

			announcementRequest := it.Event
			announcementRequest.ProviderId = ack.ProviderId
			announcementRequest.HandshakeId = ack.HandshakeId

			msg, err = tooleventmodule.WriteV0EventString(&announcementRequest)
			if err != nil {
				t.Fatal("internal test error: unable to encode event: ", it.Event)
			}

			ws.WriteMessage(websocket.TextMessage, []byte(*msg))
			ws.SetReadDeadline(time.Now().Add(time.Second))
			_, buff, err = ws.ReadMessage()
			if err != nil {
				fmt.Println("read error: ", err.Error())
				if !it.Success {
					return
				}

				t.Fatal("error response from annoucement:", err.Error())
				return
			}

			str = string(buff)
			response, err := tooleventmodule.ParseEventString(&str)
			if err != nil {
				t.Error("internal test error: malformed response event", err.Error())
				return
			}

			if !response.IsValid() {
				t.Error("respnose is not valid: ", response)
				return
			}

			if !it.ACK && response.Type != tooleventmodule.EventTypeAnnouncementNACK {
				t.Error("announcement should NOT have been acknowledged")
				return
			}

			if it.ACK && response.Type != tooleventmodule.EventTypeAnnouncementACK {
				t.Error("announcement should have been acknowledged")
				return
			}

			if !it.ACK {
				return
			}

			if response.HandlerId == "" ||
				response.ProviderId != announcementRequest.ProviderId ||
				response.HandshakeId != announcementRequest.HandshakeId {
				t.Error("incorrect response parameters:", response)
				return
			}

			if len(provider.handlers) == 0 {
				t.Error("expected to have handlers inside provider:")
				return
			}

			handler := provider.handlers[len(provider.handlers)-1]
			if handler.Tool.Acronym != it.Event.Tool {
				t.Error("most recent handler tool does not match. expected:", it.Event.Tool, "got:", handler.Tool)
				return
			}

			fmt.Println("proceeding with a command to distribute")

			command := exampleEventToDistribute
			command.ProviderId = response.ProviderId
			command.HandshakeId = response.HandshakeId
			command.HandlerId = response.HandlerId

			msg, err = tooleventmodule.WriteV0EventString(&command)
			if err != nil {
				t.Fatal("internal test error: unable to encode event: ", it.Event)
			}

			ws.WriteMessage(websocket.TextMessage, []byte(*msg))

			// just enough time to be processed
			time.Sleep(100 * time.Millisecond)

			distrubutedCount := len(manager.DistributedEvents)
			if distrubutedCount == 0 {
				t.Error("no event was distributed until now, but expected command have been")
				return
			}

			lastDistributed := manager.DistributedEvents[distrubutedCount-1]
			if !cmp.Equal(*lastDistributed, command) {
				t.Error("last distributed event differs from sent event", cmp.Diff(*lastDistributed, command))
			}
		})
	}

}

type testManager struct {
	DistributedEvents []*tooleventmodule.ToolEvent
}

// OnProviderDisconnection implements Manager.
func (t *testManager) OnProviderDisconnection(prov *Provider) {
	panic("unimplemented")
}

// DistributeEvent implements Manager.
func (t *testManager) DistributeEvent(e *tooleventmodule.ToolEvent) {
	t.DistributedEvents = append(t.DistributedEvents, e)
}

// RegisterProviderProjectAndTool implements Manager.
func (t *testManager) RegisterProviderProjectAndTool(m *ProviderToolMapping) {
	log.Printf("[testManager] should be registering project and tool (%s, %s)", m.Project.Acronym, m.Tool.Acronym)
}

// FindProject implements Manager.
func (t *testManager) FindProject(acronym string) (*projectmodule.Project, error) {
	if acronym == "validproject" {
		return &projectmodule.Project{
			ID:          0,
			Acronym:     "validproject",
			Name:        "valid project",
			Description: "valid project",
		}, nil
	}

	return nil, &commonmodule.GenericLogicError{Message: "not found"}
}

// FindTool implements Manager.
func (t *testManager) FindTool(project *projectmodule.Project, acronym string) (*toolmodule.Tool, error) {
	if acronym == "validtool" {
		return &toolmodule.Tool{
			ID:          0,
			ProjectId:   0,
			Acronym:     "validtool",
			Name:        "valid tool",
			Description: "valid tool",
		}, nil
	}

	return nil, &commonmodule.GenericLogicError{Message: "not found"}
}
