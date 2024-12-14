package api

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	modules "platformlab/controlpanel/modules/user"
	"platformlab/controlpanel/server/api/connectionmgr"
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

	userService modules.UserService

	providerMgr *connectionmgr.ProviderMgr
	clientMgr   *connectionmgr.ClientMgr
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

		userService: modules.UserService{Db: db},

		clientMgr:   &connectionmgr.ClientMgr{},
		providerMgr: nil,
	}
}
