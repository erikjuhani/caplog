package main

import (
	"log"
	"os"

	"github.com/erikjuhani/caplog/cli"
	"github.com/erikjuhani/caplog/config"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	if err := config.Load(homeDir); err != nil {
		log.Fatal(err)
	}

	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}
}
