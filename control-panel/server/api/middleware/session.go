package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	commonmodule "synthreon/modules/common"
	sessionmodule "synthreon/modules/session"
	"synthreon/server/api"
	"strings"

	"github.com/gorilla/mux"
)

// This is safe because the parameter as a variable is registered when the
// route is being registered, so we can know that it wont be passed in a route
// that shouldn't support it.
// Otherwise, if it was a query parameter this would be a risk
func tryGetPathToken(r *http.Request) *string {
	params := mux.Vars(r)
	v, ok := params["accessToken"]
	if !ok {
		return nil
	}

	return &v
}

func GetSessionMiddleware(secret string) func(next http.Handler) http.Handler {
	sessionService := sessionmodule.NewSessionService(secret)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				var accessToken string

				log.Println("[sessionMiddleware] Executing session middleware")

				authorizationHeader := r.Header.Get("Authorization")
				maybeTokenParameter := tryGetPathToken(r)
				if authorizationHeader != "" {
					authorizationHeaderParts := strings.Split(authorizationHeader, " ")
					if len(authorizationHeaderParts) < 2 || !strings.EqualFold(authorizationHeaderParts[0], "bearer") {
						log.Print("[sessionMiddleware] unexpected authorization header structure")
						w.WriteHeader(http.StatusUnauthorized)
						json.NewEncoder(w).Encode(api.ErrorMessage{Message: "session.required"})
						return
					}

					accessToken = authorizationHeaderParts[1]
				} else if maybeTokenParameter != nil {
					accessToken = *maybeTokenParameter
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
