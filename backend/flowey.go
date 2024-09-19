package main

import (
	"errors"
	"log"
	"os"
)

func run() error {
	if len(os.Args) < 2 {
		return errors.New("no address specified")
	}
	address := os.Args[1]
	server := NewFloweyServer(address)
	return server.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
