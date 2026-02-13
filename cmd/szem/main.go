package main

import (
	"context"
	"log"
	"os"

	"github.com/weihaoli/szem/pkg/repo/sqlite"
	"github.com/weihaoli/szem/pkg/server"
	"github.com/weihaoli/szem/pkg/storage"
	"github.com/weihaoli/szem/pkg/storage/local"
	"github.com/weihaoli/szem/pkg/storage/pdp"
)

func main() {

	backend := os.Getenv("QAVE_STORAGE_BACKEND")

	var st storage.Storage
	switch backend {
	case "local":
		st = local.New("./data")
	case "pdp":
		endpoint := os.Getenv("PDP_ENDPOINT")
		apiKey := os.Getenv("PDP_API_KEY")
		if endpoint == "" {
			log.Fatal("PDP_ENDPOINT not set")
		}
		st = pdp.New(endpoint, apiKey)
	default:
		log.Println("QAVE_STORAGE_BACKEND not set, using local")
		st = local.New("./data")
	}

	// SQLite
	db, err := sqlite.Open("./data/qave.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	objects := sqlite.NewObjectRepo(db)
	jobs := sqlite.NewJobRepo(db)

	ctx := context.Background()
	if err := objects.Init(ctx); err != nil {
		log.Fatal(err)
	}
	if err := jobs.Init(ctx); err != nil {
		log.Fatal(err)
	}

	srv := server.New(":8080", st, objects, jobs)
	log.Fatal(srv.Start())
}

