package server

import (
	"log"
	"net/http"

	"github.com/weihaoli/szem/pkg/s3"
)

type Server struct {
	addr string
}

func New(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Start() error {
	handler := s3.NewIngress()

	handler = WithRequestID(handler)

	log.Printf("SZEM listening on %s", s.addr)

	return http.ListenAndServe(s.addr, handler)
}


