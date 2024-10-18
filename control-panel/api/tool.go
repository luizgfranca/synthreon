package api

import (
	"encoding/json"
	"log"
	"net/http"
	"platformlab/controlpanel/api/connectionmgr"
	"platformlab/controlpanel/model"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Tool struct {
	websocketUpgrader websocket.Upgrader

	providerMgr *connectionmgr.ProviderMgr
	clientMgr   *connectionmgr.ClientMgr
}

func (t *Tool) GetEventRresponseTEST() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input model.ToolEvent
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			json.NewEncoder(w).Encode(ErrorMessage{Message: err.Error()})
			return
		}

		if !input.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMessage{Message: "invalid request data"})
			return
		}

		helloWorldMock := model.ToolEvent{
			Class:   model.EventClassOperation,
			Type:    model.EventTypeDisplay,
			Project: "x",
			Tool:    "y",
			Display: &model.DisplayDefniition{
				Type: model.DisplayDefniitionTypeView,
				Elements: &[]model.DisplayElement{
					{Type: "heading", Text: "Hello world", Description: "This is a hello world test"},
				},
			},
		}

		json.NewEncoder(w).Encode(helloWorldMock)
	}
}

func (t *Tool) clientConnectionHandler(connection *websocket.Conn) {
	log.Print("[tool] new client connection")
	t.clientMgr.NewClient(connection)
}

func (t *Tool) providerConnectionHandler(connection *websocket.Conn) {
	t.providerMgr = connectionmgr.NewProviderConnectionMgr(connection)
}

func upgraderAllowAllOrigins(r *http.Request) bool {
	return true
}

func (t *Tool) ToolClientWebsocket() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		connection, err := t.websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("websocket upgrade error: ", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMessage{Message: err.Error()})
			return
		}

		go t.clientConnectionHandler(connection)
	}
}

func (t *Tool) ToolProviderWebsocket() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		connection, err := t.websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("websocket upgrade error: ", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorMessage{Message: err.Error()})
			return
		}

		go t.providerConnectionHandler(connection)
	}
}

func ToolRestAPI(db *gorm.DB) *Tool {
	return &Tool{websocket.Upgrader{CheckOrigin: upgraderAllowAllOrigins}, nil, &connectionmgr.ClientMgr{}}
}
