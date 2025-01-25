package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	commonmodule "platformlab/controlpanel/modules/common"
	sessionmodule "platformlab/controlpanel/modules/session"
	"platformlab/controlpanel/server/api"
	"strings"
)

func GetSessionMiddleware(secret string) func(next http.Handler) http.Handler {
	sessionService := sessionmodule.NewSessionService(secret)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				var accessToken string

				log.Println("[sessionMiddleware] Executing session middleware")

				authorizationHeader := r.Header.Get("Authorization")
				secureWebsocketProtocolHeader := r.Header.Get("Sec-WebSocket-Protocol")
				if authorizationHeader != "" {
					authorizationHeaderParts := strings.Split(authorizationHeader, " ")
					if len(authorizationHeaderParts) < 2 || !strings.EqualFold(authorizationHeaderParts[0], "bearer") {
						log.Print("[sessionMiddleware] unexpected authorization header structure")
						w.WriteHeader(http.StatusUnauthorized)
						json.NewEncoder(w).Encode(api.ErrorMessage{Message: "session.required"})
						return
					}

					accessToken = authorizationHeaderParts[1]
				} else if secureWebsocketProtocolHeader != "" {
					accessToken = secureWebsocketProtocolHeader
				} else {
					log.Print("[sessionMiddleware] unable to get authorization header")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(api.ErrorMessage{Message: "session.required"})
					return
				}

				session, err := sessionService.DecodeToken(accessToken)
				if err != nil {
					log.Print("[sessionMiddleware] error loading session token: ", err.Error())
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(api.ErrorMessage{Message: err.Error()})
					return
				}

				ctx := context.WithValue(r.Context(), commonmodule.SessionRequestContextKey, session)

				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}
