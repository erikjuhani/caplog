package main

import (
	"log"

	"github.com/erikjuhani/caplog/cmd"
	"github.com/erikjuhani/caplog/config"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}
