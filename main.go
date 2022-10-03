package main

import (
	"log"

	"github.com/erikjuhani/caplog/cli"
	"github.com/erikjuhani/caplog/config"
)

func main() {
	log.SetFlags(0)

	if err := config.Load(); err != nil {
		log.Fatalf("error: %s", err)
	}

	if err := cli.Run(); err != nil {
		log.Fatalf("error: %s", err)
	}
}
