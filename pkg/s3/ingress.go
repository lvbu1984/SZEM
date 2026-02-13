package s3

import (
	"net/http"

	"github.com/weihaoli/szem/pkg/repo"
	"github.com/weihaoli/szem/pkg/storage"
)

type Ingress struct {
	storage storage.Storage
	objects repo.ObjectRepo
	jobs    repo.JobRepo
}

func NewIngress(s storage.Storage, objects repo.ObjectRepo, jobs repo.JobRepo) http.Handler {
	return &Ingress{
		storage: s,
		objects: objects,
		jobs:    jobs,
	}
}

func (i *Ingress) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Handle(w, r, i.storage, i.objects, i.jobs)
}

