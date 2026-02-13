package server

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestID := uuid.New().String()

		w.Header().Set("X-Request-Id", requestID)

		log.Printf("[REQ %s] %s %s", requestID, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

