package toolentity_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"synthreon/modules/toolentity"
	tooleventmodule "synthreon/modules/toolevent"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/websocket"
)

func upgraderAllowAllOrigins(r *http.Request) bool {
	return true
}

// NOTICE:
//
// this is a VERY basic test, it just validates if the
// base logic is working, but doesn't try to really break
// the logic yet
func TestToolEventRules(t *testing.T) {
	fmt.Println("test WebSocket tool entity")

	testSendEvent := tooleventmodule.ToolEvent{
		Type:        tooleventmodule.EventTypeHandshakeRequest,
		Project:     "x",
		HandshakeId: "hsh",
	}
	testResponseEvent := tooleventmodule.ToolEvent{
		Type:        tooleventmodule.EventTypeHandshakeACK,
		Project:     "x",
		HandshakeId: "hsh",
		ProviderId:  "prov",
	}

	if !testSendEvent.IsValid() || !testResponseEvent.IsValid() {
		t.Fatal("internal test error: send or response events are not valid")
	}

	msgToSend, err := tooleventmodule.WriteV0EventString(&testSendEvent)
	if err != nil {
		t.Fatal("internal test error: unable to serialize event: ", err.Error())
	}

	var receivedEventRegister *tooleventmodule.ToolEvent
	var errorResponding error

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		websocketUpgrader := websocket.Upgrader{CheckOrigin: upgraderAllowAllOrigins}
		conn, err := websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal("internal test error: ", err.Error())
		}

		entity := toolentity.NewWebsocketToolEntity(conn)
		if entity == nil {
			t.Error("null websocket entity after creation")
		}

		entity.OnEventReceived(func(event *tooleventmodule.ToolEvent) {
			fmt.Println("[test] OnEventReceived: ", *event)
			receivedEventRegister = event
			entity.SendEvent(&testResponseEvent)
			if err != nil {
				errorResponding = err
			}
		})

		entity.StartHandler()
	}))
	defer s.Close()

	// TODO: should test sending
	// - empty message
	// - malformed message
	// - event with invalid logic
	// - multiple events
	// - invlaid events in the middle of valid events
	// the same should be done for responses
	websocketURL := "ws" + strings.TrimPrefix(s.URL, "http")
	fmt.Println("websocket URL: ", websocketURL)
	connection, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		t.Fatal("internal test error: unable to connect to test websocket: ", err.Error())
	}

	err = connection.WriteMessage(websocket.TextMessage, []byte(*msgToSend))
	if err != nil {
		t.Fatal("intenal test erro: unable to send message to websocket server: ", err.Error())
	}

	// just relinquishing control to let processing execute
	time.Sleep(time.Millisecond)

	if receivedEventRegister == nil {
		t.Error("testSendEvent was not received by handler")
		return
	}

	if !cmp.Equal(*receivedEventRegister, testSendEvent) {
		t.Error(
			"received event is different from the expected\n Received:", *receivedEventRegister,
			"\n Expected:", testSendEvent,
			"\n", cmp.Diff(*receivedEventRegister, testSendEvent),
		)

		return
	}

	if errorResponding != nil {
		t.Error("error responding to event: ", errorResponding.Error())
		return
	}

	_, msg, err := connection.ReadMessage()
	if err != nil {
		t.Fatal("intenal test erro: error reading response message: ", err.Error())
	}

	msgStr := string(msg)
	event, err := tooleventmodule.ParseEventString(&msgStr)
	if err != nil {
		t.Error("unable to parse response event: ", event)
		return
	}

	if !cmp.Equal(*event, testResponseEvent) {
		t.Error(
			"received event is different from the expected\n Received:", *event,
			"\n Expected:", testResponseEvent,
			"\n", cmp.Diff(*event, testResponseEvent),
		)

		return
	}
}
