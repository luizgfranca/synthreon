package providermodule

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	projectmodule "platformlab/controlpanel/modules/project"
	toolmodule "platformlab/controlpanel/modules/tool"
	"platformlab/controlpanel/modules/toolentity"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
	tooleventdisplay "platformlab/controlpanel/modules/toolevent/display"
	tooleventinput "platformlab/controlpanel/modules/toolevent/input"
	tooleventresult "platformlab/controlpanel/modules/toolevent/result"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/websocket"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var lastForwardedEvent *tooleventmodule.ToolEvent

type testOrchestrator struct{}

// ForwardEvent implements Orchestrator.
// TODO: tests here are only the essential to advance with development,
// should improve them soon
func TestBasicProviderManagerBehavior(t *testing.T) {

	fmt.Println("starting in-mmemory database")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed connecting to database")
	}

	projectService := projectmodule.ProjectService{Db: db}
	toolService := toolmodule.ToolService{Db: db}

	fmt.Println("doing migrations")
	db.AutoMigrate(&projectmodule.Project{})
	db.AutoMigrate(&toolmodule.Tool{})

	fmt.Println("creating mock data")

	exampleProject := &projectmodule.Project{
		ID:          1,
		Acronym:     "proj",
		Name:        "Project",
		Description: "Project",
	}
	log.Println(exampleProject)
	projectService.Create(exampleProject)

	exampleToolA := toolmodule.Tool{
		ID:          1,
		Acronym:     "ta",
		Name:        "Tool A",
		Description: "Tool A",
		ProjectId:   1,
	}
	log.Println(exampleToolA)
	toolService.Create(&exampleToolA)

	exampleToolB := toolmodule.Tool{
		ID:          2,
		Acronym:     "tb",
		Name:        "Tool B",
		Description: "Tool B",
		ProjectId:   1,
	}
	log.Println(exampleToolB)
	toolService.Create(&exampleToolB)

	tools := toolService.FindAll()
	log.Println(tools)

	orchestrator := testOrchestrator{}
	providerManager := NewProviderManagerService(
		orchestrator,
		&projectService,
		&toolService,
	)

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

		providerManager.EntityConnection(entity)
	}))
	defer s.Close()

	websocketURL := "ws" + strings.TrimPrefix(s.URL, "http")
	fmt.Println("websocket URL: ", websocketURL)
	ws, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatal("internal test error: unable to connect to test websocket: ", err.Error())
	}

	msg, err := tooleventmodule.WriteV0EventString(&tooleventmodule.ToolEvent{
		Type:    tooleventmodule.EventTypeHandshakeRequest,
		Project: "proj",
	})
	if err != nil {
		t.Fatal("internal test error: unable to encode event")
	}

	ws.WriteMessage(websocket.TextMessage, []byte(*msg))
	ws.SetReadDeadline(time.Now().Add(time.Second))
	_, buff, err := ws.ReadMessage()
	if err != nil {
		t.Error("ack message not received provider", err.Error())
		return
	}

	str := string(buff)
	event, err := tooleventmodule.ParseEventString(&str)
	if err != nil {
		t.Error("malformed response event", err.Error())
		return
	}

	if event.Type != tooleventmodule.EventTypeHandshakeACK {
		t.Error("handshake should have been acknowleged, response was: ", event)
		return
	}

	if event.Project != "proj" || event.HandshakeId == "" || event.ProviderId == "" {
		t.Error("invalid response, response was: ", event)
		return
	}

	var provider *Provider
	for _, p := range providerManager.providers {
		if event.ProviderId == p.ID {
			provider = p
			break
		}
	}

	if provider == nil {
		t.Error("no provider found inside providermanager for provider retorned in the event: ", event)
		return
	}

	if provider.handshakeId != event.HandshakeId {
		t.Error("provider's registered handshake does not match respnse's: ", event)
		return
	}

	log.Println("projects", projectService.FindAll())

	if !cmp.Equal(*provider.Project, *exampleProject) {
		t.Error("project does not match expected: \n", cmp.Diff(*provider.Project, exampleProject))
		return
	}

	msg, err = tooleventmodule.WriteV0EventString(&tooleventmodule.ToolEvent{
		Type:        tooleventmodule.EventTypeAnnouncementHandler,
		Project:     "proj",
		Tool:        "ta",
		HandshakeId: event.HandshakeId, // to be fiiled
		ProviderId:  event.ProviderId,  // to be fiiled
	})
	if err != nil {
		t.Fatal("internal test error: unable to encode event")
	}

	ws.WriteMessage(websocket.TextMessage, []byte(*msg))
	ws.SetReadDeadline(time.Now().Add(time.Second))
	_, buff, err = ws.ReadMessage()
	if err != nil {
		t.Error("announcement received by provider", err.Error())
		return
	}

	str = string(buff)
	event, err = tooleventmodule.ParseEventString(&str)
	if err != nil {
		t.Error("malformed response event", err.Error())
		return
	}

	if event.Type != tooleventmodule.EventTypeAnnouncementACK {
		t.Error("handshake should have been acknowleged, response was: ", event)
		return
	}

	if event.Project != "proj" && event.Tool != "ta" {
		t.Error("handshake project and tool do not match requested, response was: ", event)
		return
	}

	if event.ProviderId != provider.ID {
		t.Error("provider is not the same as handshake one, expected: ", provider.ID, "got: ", event.ProviderId)
		return
	}

	exampleOpenEvent := &tooleventmodule.ToolEvent{
		Type:      tooleventmodule.EventTypeInteractionOpen,
		Project:   "proj",
		Tool:      "ta",
		ContextId: "ctx",
	}
	providerManager.SendEvent(exampleOpenEvent)

	ws.SetReadDeadline(time.Now().Add(time.Second))
	_, buff, err = ws.ReadMessage()
	if err != nil {
		t.Error("open message not sent to provider", err.Error())
		return
	}

	str = string(buff)
	event, err = tooleventmodule.ParseEventString(&str)
	if err != nil {
		t.Error("malformed response event", err.Error())
		return
	}

	if event.Type != exampleOpenEvent.Type ||
		event.Project != exampleOpenEvent.Project ||
		event.Tool != exampleOpenEvent.Tool ||
		event.ContextId != exampleOpenEvent.ContextId ||
		event.ProviderId != provider.ID {
		t.Error(
			"event sent to provider should match sent information\n",
			"got: ", event, "\n",
			"original event:", exampleOpenEvent, "\n",
			"provider: ", provider, "\n",
		)
	}

	exampleDisplayCommand := &tooleventmodule.ToolEvent{
		Type:        tooleventmodule.EventTypeCommandDisplay,
		Project:     event.Project,
		Tool:        event.Tool,
		ContextId:   event.ContextId,
		ProviderId:  event.ProviderId,
		HandlerId:   event.HandlerId,
		ExecutionId: "exec",
		Display: &tooleventdisplay.DisplayDefinition{
			TextBox: &tooleventdisplay.TextBoxDisplay{
				Content: "test",
			},
		},
	}

	msg, err = tooleventmodule.WriteV0EventString(exampleDisplayCommand)
	if err != nil {
		t.Fatal("internal test error: unable to encode event")
	}

	ws.WriteMessage(websocket.TextMessage, []byte(*msg))
	time.Sleep(100 * time.Millisecond)

	if lastForwardedEvent == nil {
		t.Error("expected display event to have been forwarded")
		return
	}

	if lastForwardedEvent.Type != tooleventmodule.EventTypeCommandDisplay ||
		lastForwardedEvent.ContextId != "ctx" ||
		!cmp.Equal(*lastForwardedEvent.Display, *exampleDisplayCommand.Display) {
		t.Error(
			"last forwarded event does not match test display command sent\n",
			"sent: ", exampleDisplayCommand, "\n",
			"last forwarded", *lastForwardedEvent, "\n",
			"display diff:", cmp.Diff(*lastForwardedEvent.Display, *exampleDisplayCommand.Display), "\n",
		)
		return
	}

	exampleOKInput := tooleventmodule.ToolEvent{
		Type:       tooleventmodule.EventTypeInteractionInput,
		Project:    "proj",
		Tool:       "ta",
		ContextId:  "ctx",
		TerminalId: "tid",
		Input:      &tooleventinput.InputDefinition{},
	}

	err = providerManager.SendEvent(&exampleOKInput)
	if err != nil {
		t.Fatal("internal test eror: error sendint test input: ", exampleOKInput, err.Error())
	}

	ws.SetReadDeadline(time.Now().Add(time.Second))
	_, buff, err = ws.ReadMessage()
	if err != nil {
		t.Error("input interaction not sent to provider", err.Error())
		return
	}

	str = string(buff)
	event, err = tooleventmodule.ParseEventString(&str)
	if err != nil {
		t.Error("malformed response event", err.Error())
		return
	}

	if event.Type != tooleventmodule.EventTypeInteractionInput ||
		event.ContextId != "ctx" {
		t.Error("event received does not match sent", event)
		return
	}

	if event.ExecutionId != "exec" {
		t.Error("executionId was not filled", event)
		return
	}

	exampleOKInput.Tool = "tb"

	err = providerManager.SendEvent(&exampleOKInput)
	if err == nil {
		t.Error("should not have succeeded in sending event to unregistered tool")
	}

	exampleFinishCommand := &tooleventmodule.ToolEvent{
		Type:        tooleventmodule.EventTypeCommandFinish,
		Project:     event.Project,
		Tool:        event.Tool,
		ContextId:   event.ContextId,
		ProviderId:  event.ProviderId,
		HandlerId:   event.HandlerId,
		ExecutionId: "exec",
		Result: &tooleventresult.ToolEventResult{
			Status:  tooleventresult.ToolEventResultStatusSuccess,
			Message: "success",
		},
	}

	msg, err = tooleventmodule.WriteV0EventString(exampleFinishCommand)
	if err != nil {
		t.Fatal("internal test error: unable to encode event")
	}

	ws.WriteMessage(websocket.TextMessage, []byte(*msg))
	time.Sleep(100 * time.Millisecond)

	if lastForwardedEvent == nil {
		t.Error("expected finish event to have been forwarded")
		return
	}

	if lastForwardedEvent.Type != tooleventmodule.EventTypeCommandFinish {
		t.Error("last forwarded event does not match a finish command")
	}
}

func (t testOrchestrator) ForwardEvent(e *tooleventmodule.ToolEvent) {
	lastForwardedEvent = e
}
