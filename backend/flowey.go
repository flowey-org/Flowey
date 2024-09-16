package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
)

func run() error {
	if len(os.Args) < 2 {
		return errors.New("no address specified")
	}
	address := os.Args[1]
	server := &http.Server{Addr: address, Handler: Handler{}}
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	log.Printf("listening at %v", address)
	return server.Serve(listener)
}

func main() {
	log.Fatal(run())
}
