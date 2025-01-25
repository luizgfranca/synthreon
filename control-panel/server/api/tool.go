package api

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	commonmodule "platformlab/controlpanel/modules/common"
	orchestratormodule "platformlab/controlpanel/modules/orchestrator"
	projectmodule "platformlab/controlpanel/modules/project"
	sessionmodule "platformlab/controlpanel/modules/session"
	toolmodule "platformlab/controlpanel/modules/tool"
	"platformlab/controlpanel/modules/toolentity"
	modules "platformlab/controlpanel/modules/user"
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

	userService         modules.UserService
	orchestratorService *orchestratormodule.OrchestratorService
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

		// TODO: abstract this behavior later
		v := r.Context().Value(commonmodule.SessionRequestContextKey)
		session, ok := v.(*sessionmodule.Session)
		if !ok {
			log.Fatalln("session expected in context when starting websocket connection")
		}

		user, err := t.userService.FindByEmail(session.Email)
		if err != nil {
			log.Fatalln("unexpected: session ", session, " emitted with an invalid user?")
		}

		entity := toolentity.NewWebsocketToolEntity(connection)
		t.orchestratorService.RegisterClientEntity(session, user, entity)
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

		entity := toolentity.NewWebsocketToolEntity(connection)
		t.orchestratorService.RegisterProviderEntity(entity)
	}
}

func ToolRestAPI(db *gorm.DB) *Tool {
	toolService := toolmodule.ToolService{Db: db}
	projectService := projectmodule.ProjectService{Db: db}

	return &Tool{
		websocketUpgrader: websocket.Upgrader{CheckOrigin: upgraderAllowAllOrigins},

		userService: modules.UserService{Db: db},
		orchestratorService: orchestratormodule.NewOrchestratorService(
			&projectService,
			&toolService,
		),
	}
}
