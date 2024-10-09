package main

import (
	"log"
)

func run() error {
	address := "0.0.0.0:80"
	server := NewServer(address)
	return server.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
