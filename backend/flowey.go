package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
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

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	errs := make(chan error, 1)
	go func() {
		errs <- server.Serve(listener)
	}()

	select {
	case <-sigint:
		fmt.Print("\rinterrupting...\n")
	case err := <-errs:
		log.Printf("failed to serve: %w", err)
	}

	return server.Shutdown(context.Background())
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
