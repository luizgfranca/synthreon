package middleware

import (
	"net/http"
)

func GetCORSMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			/**
			 * TODO: there should be a specific module to generate messages,
			 *		 it should used in a middleware to generate a requestId for each request and add it
			 *		 to a context, and the log service would add this ID to all messages related to that
			 *		 request.
			 *
			 *		 Also, there should be a middleware to add this requestId in the response header
			 */

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			next.ServeHTTP(w, r)
		})
	}
}
