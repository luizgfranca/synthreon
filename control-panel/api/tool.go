package api

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"platformlab/controlpanel/api/connectionmgr"
	"platformlab/controlpanel/model"
	"platformlab/controlpanel/service"
	"strings"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

const (
	CredentialsUsernamePosition = 0
	CredentialsPasswordPosition = 1
)

type Tool struct {
	websocketUpgrader websocket.Upgrader

	userService service.User

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
	log.Print("[ToolAPI] new client connection")
	t.clientMgr.NewClient(connection)
}

func (t *Tool) providerConnectionHandler(connection *websocket.Conn) {
	log.Print("[ToolAPI] new provider connection")
	t.providerMgr = connectionmgr.NewProviderConnectionMgr(connection)

	// TODO: maybe it should be better if instead of passing them i just passed the
	// 		 complete tool context that would have the managers
	t.clientMgr.ProviderMgr = t.providerMgr
	t.providerMgr.ClientMgr = t.clientMgr
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

		log.Println("[ToolApi] validating credentials")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			log.Println("[ToolApi] expected authorization header")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorMessage{Message: "provider.unauthorized"})
			return
		}

		log.Println("[ToolApi] authorizationHeader: ", authorizationHeader)

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) < 2 || !strings.EqualFold(headerParts[0], "basic") {
			log.Println("[ToolApi] malformed authorization header")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorMessage{Message: "provider.unauthorized"})
			return
		}

		credentialsStr, err := base64.StdEncoding.DecodeString(headerParts[1])
		if err != nil {
			log.Println("[ToolApi] invalid header's credentials string")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorMessage{Message: "provider.unauthorized"})
			return
		}

		credentials := strings.Split(string(credentialsStr), ":")
		if len(credentials) < 2 {
			log.Println("[ToolApi] malformed credentials")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorMessage{Message: "provider.unauthorized"})
			return
		}

		_, err = t.userService.VerifyAuthenticationCredentials(
			&credentials[CredentialsUsernamePosition],
			&credentials[CredentialsPasswordPosition],
		)
		if err != nil {
			log.Println("[ToolApi] invalid user credentials")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorMessage{Message: "provider.unauthorized"})
			return
		}

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
	return &Tool{
		websocketUpgrader: websocket.Upgrader{CheckOrigin: upgraderAllowAllOrigins},

		userService: service.User{Db: db},

		clientMgr:   &connectionmgr.ClientMgr{},
		providerMgr: nil,
	}
}
