package clientmodule

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	projectmodule "platformlab/controlpanel/modules/project"
	sessionmodule "platformlab/controlpanel/modules/session"
	toolmodule "platformlab/controlpanel/modules/tool"
	"platformlab/controlpanel/modules/toolentity"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
	tooleventdisplay "platformlab/controlpanel/modules/toolevent/display"
	tooleventinput "platformlab/controlpanel/modules/toolevent/input"
	tooleventresult "platformlab/controlpanel/modules/toolevent/result"
	usermodule "platformlab/controlpanel/modules/user"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/websocket"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testOrchestrator struct{}

var lastForwardedEvent *tooleventmodule.ToolEvent

func TestClientManagerBaseBehavior(t *testing.T) {

	fmt.Println("starting in-mmemory database")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed connecting to database")
	}

	projectService := projectmodule.ProjectService{Db: db}
	toolService := toolmodule.ToolService{Db: db}
	userService := usermodule.UserService{Db: db}

	fmt.Println("doing migrations")
	db.AutoMigrate(&projectmodule.Project{})
	db.AutoMigrate(&toolmodule.Tool{})
	db.AutoMigrate(&usermodule.User{})

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

	testOpenEvent := tooleventmodule.ToolEvent{
		Type:      tooleventmodule.EventTypeInteractionOpen,
		Project:   "proj",
		Tool:      "ta",
		SessionId: "", // will be filled when session is created
	}

	testDisplayEvent := tooleventmodule.ToolEvent{
		Type:      tooleventmodule.EventTypeCommandDisplay,
		Project:   "proj",
		Tool:      "ta",
		ContextId: "", // to be filled
		Display: &tooleventdisplay.DisplayDefinition{
			TextBox: &tooleventdisplay.TextBoxDisplay{
				Content: "content",
			},
		},
	}

	testInputEvent := tooleventmodule.ToolEvent{
		Type:       tooleventmodule.EventTypeInteractionInput,
		Project:    "proj",
		Tool:       "ta",
		ContextId:  "", // to be filled
		TerminalId: "", // to be filled
		SessionId:  "", // to be filled
		Input: &tooleventinput.InputDefinition{
			Fields: []tooleventinput.InputField{},
		},
	}

	testFinishEvent := tooleventmodule.ToolEvent{
		Type:      tooleventmodule.EventTypeCommandFinish,
		Project:   "proj",
		Tool:      "ta",
		ContextId: "", // to be filled
		Result: &tooleventresult.ToolEventResult{
			Status:  tooleventresult.ToolEventResultStatusSuccess,
			Message: "msg",
		},
	}

	exampleUser, err := usermodule.NewUser(
		"test user",
		"test@test.com",
		"pass",
	)
	if err != nil {
		panic("error creating mock user")
	}
	userService.Create(exampleUser)

	session := sessionmodule.NewSessionFromUser(exampleUser)
	testOpenEvent.SessionId = session.ID

	orchestrator := testOrchestrator{}

	clientManager := NewCLientManagerService(
		orchestrator,
		&projectService,
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

		clientManager.EntityConnection(session, exampleUser, entity)
	}))
	defer s.Close()

	websocketURL := "ws" + strings.TrimPrefix(s.URL, "http")
	fmt.Println("websocket URL: ", websocketURL)
	ws, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatal("internal test error: unable to connect to test websocket: ", err.Error())
	}

	msg, err := tooleventmodule.WriteV0EventString(&testOpenEvent)
	if err != nil {
		t.Fatal("internal test error: unable to encode event")
	}

	err = ws.WriteMessage(websocket.TextMessage, []byte(*msg))
	if err != nil {
		panic("internal test error: unable to send message")
	}

	time.Sleep(100 * time.Millisecond)

	if lastForwardedEvent == nil {
		t.Error("expected open event to have been forwarded")
		return
	}

	log.Println("[test] validating forwarded event")
	if lastForwardedEvent.Type != tooleventmodule.EventTypeInteractionOpen ||
		lastForwardedEvent.ContextId == "" ||
		lastForwardedEvent.Tool != testOpenEvent.Tool ||
		lastForwardedEvent.Project != testOpenEvent.Project {
		t.Error(
			"last forwarded event does not match test open interaction sent\n",
			"sent: ", testOpenEvent, "\n",
			"last forwarded", *lastForwardedEvent, "\n",
		)
		return
	}

	log.Println("[test] sending display response")
	contextId := lastForwardedEvent.ContextId
	testDisplayEvent.ContextId = contextId
	clientManager.SendEvent(&testDisplayEvent)

	log.Println("[test] receiving display response")
	ws.SetReadDeadline(time.Now().Add(time.Second))
	_, buff, err := ws.ReadMessage()
	if err != nil {
		t.Error("display command not sent to client entity", err.Error())
		return
	}

	log.Println("[test] validating display response")
	str := string(buff)
	event, err := tooleventmodule.ParseEventString(&str)
	if err != nil {
		t.Error("malformed response event", err.Error())
		return
	}

	if event.Type != tooleventmodule.EventTypeCommandDisplay ||
		event.ContextId != contextId ||
		event.TerminalId == "" ||
		event.Tool != testDisplayEvent.Tool ||
		event.Project != testDisplayEvent.Project ||
		!cmp.Equal(*event.Display, *testDisplayEvent.Display) {
		t.Error(
			"received event does not match test display command sent\n",
			"sent: ", testDisplayEvent, "\n",
			"last forwarded", *event, "\n",
			"display diff:", cmp.Diff(*testDisplayEvent.Display, *event.Display), "\n",
		)
		return
	}

	log.Println("[test] sending input")
	terminalId := event.TerminalId

	testInputEvent.ContextId = contextId
	testInputEvent.TerminalId = terminalId

	msg, err = tooleventmodule.WriteV0EventString(&testInputEvent)
	if err != nil {
		t.Fatal("internal test error: unable to encode event")
	}

	err = ws.WriteMessage(websocket.TextMessage, []byte(*msg))
	if err != nil {
		t.Fatal("internal test error: unable to send message")
	}

	log.Println("[test] receiving forwarded event")
	time.Sleep(100 * time.Millisecond)

	if lastForwardedEvent == nil {
		t.Error("expected input event to have been forwarded")
		return
	}

	if lastForwardedEvent.Type != tooleventmodule.EventTypeInteractionInput ||
		lastForwardedEvent.ContextId != contextId ||
		lastForwardedEvent.Tool != testOpenEvent.Tool ||
		lastForwardedEvent.Project != testOpenEvent.Project ||
		!cmp.Equal(*lastForwardedEvent.Input, *testInputEvent.Input) {
		{
			t.Error(
				"last forwarded event does not match test input interaction sent\n",
				"sent: ", testOpenEvent, "\n",
				"last forwarded", *lastForwardedEvent, "\n",
				"input diff:", cmp.Diff(*lastForwardedEvent.Input, *testInputEvent.Input), "\n",
			)
			return
		}
	}

	log.Println("[test] sending finish event")
	testFinishEvent.ContextId = contextId

	err = clientManager.SendEvent(&testFinishEvent)
	if err != nil {
		t.Fatal("internal test error: unable to encode event")
	}

	log.Println("[test] receiving finish event")
	ws.SetReadDeadline(time.Now().Add(time.Second))
	_, buff, err = ws.ReadMessage()
	if err != nil {
		t.Error("finish command not sent to client entity", err.Error())
		return
	}

	str = string(buff)
	event, err = tooleventmodule.ParseEventString(&str)
	if err != nil {
		t.Error("malformed response event", err.Error())
		return
	}

	log.Println("[test] validating finish event")

	if event.Type != tooleventmodule.EventTypeCommandFinish ||
		event.ContextId != contextId ||
		event.Tool != testOpenEvent.Tool ||
		event.Project != testOpenEvent.Project ||
		!cmp.Equal(*event.Result, *testFinishEvent.Result) {
		t.Error(
			"received event does not match test finish command sent\n",
			"sent: ", testFinishEvent, "\n",
			"last forwarded", *event, "\n",
			"display diff:", cmp.Diff(*event.Result, *testFinishEvent.Result), "\n",
		)
		return
	}

	clientManager.UnregisterContext(contextId)

	if clientManager.contextClientResolver.Resolve(contextId) != nil {
		t.Error("context should have been removed after finish command")
		return
	}
}

// ForwardEventToProvider implements Orchestrator.
func (t testOrchestrator) ForwardEventToProvider(e *tooleventmodule.ToolEvent) {
	lastForwardedEvent = e
}
