package middleware

import (
	"net/http"
)

type ContentType string

const (
	ContentTypeJSON         ContentType = "application/json"
	ContentTypeHTML         ContentType = "text/html"
	ContentTypeNotSpecified ContentType = ""
)

func GetContentTypeMiddleware(ct ContentType) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ct != ContentTypeNotSpecified {
				w.Header().Set("Content-Type", string(ct))
			}
			next.ServeHTTP(w, r)
		})
	}
}
