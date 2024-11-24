package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"platformlab/controlpanel/api"
	"platformlab/controlpanel/service"
	"strings"
)

func SessionMiddleware(secret string) func(next http.Handler) http.Handler {
	sessionService := service.NewSessionService(secret)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				log.Println("[sessionMiddleware] Executing session middleware")

				authorizationHeader := r.Header.Get("Authorization")
				if authorizationHeader == "" {
					log.Print("[sessionMiddleware] unable to get authorization header")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(api.ErrorMessage{Message: "session.required"})
					return
				}

				authorizationHeaderParts := strings.Split(authorizationHeader, " ")
				if len(authorizationHeaderParts) < 2 || !strings.EqualFold(authorizationHeaderParts[0], "bearer") {
					log.Print("[sessionMiddleware] unexpected authorization header structure")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(api.ErrorMessage{Message: "session.required"})
					return
				}

				accessToken := authorizationHeaderParts[1]

				session, err := sessionService.DecodeToken(accessToken)
				if err != nil {
					log.Print("[sessionMiddleware] error loading session token: ", err.Error())
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(api.ErrorMessage{Message: err.Error()})
					return
				}

				ctx := context.WithValue(r.Context(), RequestContextSession, session)

				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}
