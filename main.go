package main

import (
	"log"

	"github.com/erikjuhani/caplog/cli"
	"github.com/erikjuhani/caplog/config"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}
}
