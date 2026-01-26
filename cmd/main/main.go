package main

import (
	"log"
	"os"

	"github.com/wavy-cat/doki/cmd"
)

var (
	execute = cmd.Execute
	exit    = os.Exit
)

func run() error {
	return execute()
}

func main() {
	if err := run(); err != nil {
		log.Print(err)
		exit(1)
	}
}
