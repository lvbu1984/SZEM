package s3

import (
	"net/http"
)

// ---- SZEM Ingress entrypoint ----

type Ingress struct{}

func NewIngress() http.Handler {
	return &Ingress{}
}

func (i *Ingress) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// 所有请求统一交给 Handle
	Handle(w, r)
}

