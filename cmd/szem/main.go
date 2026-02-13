package main

import (
	"log"

	"github.com/weihaoli/szem/pkg/server"
)

func main() {
	s := server.New(":8080")

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}

