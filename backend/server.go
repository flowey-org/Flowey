package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

type FloweyServer struct {
	http.Server
}

func NewFloweyServer(address string) *FloweyServer {
	var server FloweyServer
	server.Addr = address
	server.Handler = &handler{}
	return &server
}

func (server *FloweyServer) ListenAndServe() error {
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}
	log.Printf("listening at %v", server.Addr)

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
