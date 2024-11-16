package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"platformlab/controlpanel/api"
	"platformlab/controlpanel/service"
)

func SessionMiddleware(secret string) func(next http.Handler) http.Handler {
	sessionService := service.NewSessionService(secret)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				log.Println("Executing authentication middleware")

				cookie, err := r.Cookie("session")
				if err != nil {
					log.Print("error getting session token: ", err.Error())
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(api.ErrorMessage{Message: err.Error()})
					return
				}

				if cookie == nil {
					log.Print("session cookit not found on request")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(api.ErrorMessage{Message: "session.required"})
					return
				}

				session, err := sessionService.DecodeToken(cookie.Value)
				if err != nil {
					log.Print("error loading session token: ", err.Error())
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(api.ErrorMessage{Message: err.Error()})
				}

				ctx := context.WithValue(r.Context(), RequestContextSession, session)

				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}
