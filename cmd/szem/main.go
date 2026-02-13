package main

import (
	"log"
	"net/http"
	"os"

	"github.com/your-org/szem/internal/s3ingress"
)

func main() {
	addr := os.Getenv("SZEM_LISTEN")
	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()
	mux.Handle("/", s3ingress.NewHandler())

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("SZEM listening on %s\n", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

