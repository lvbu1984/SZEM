package server

import (
	"log"
	"net/http"

	"github.com/weihaoli/szem/pkg/core"
	"github.com/weihaoli/szem/pkg/repo"
	"github.com/weihaoli/szem/pkg/s3"
	"github.com/weihaoli/szem/pkg/storage"
)

type Server struct {
	addr    string
	storage storage.Storage
	objects repo.ObjectRepo
	jobs    repo.JobRepo
}

func New(addr string, st storage.Storage, objects repo.ObjectRepo, jobs repo.JobRepo) *Server {
	return &Server{
		addr:    addr,
		storage: st,
		objects: objects,
		jobs:    jobs,
	}
}

func (s *Server) Start() error {
	// 启动 worker（最小B）
	w := &core.Worker{
		Jobs:    s.jobs,
		Objects: s.objects,
		Storage: s.storage,
	}
	w.Start()

	handler := s3.NewIngress(s.storage, s.objects, s.jobs)

	log.Printf("SZEM listening on %s", s.addr)
	return http.ListenAndServe(s.addr, handler)
}

