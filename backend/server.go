package main

import (
	"context"
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
		print("\r")
	case err := <-errs:
		log.Printf("failed to serve: %v", err)
	}

	return server.Shutdown(context.Background())
}

func (server *FloweyServer) Shutdown(ctx context.Context) error {
	log.Println("shutdown initiated")
	defer log.Println("shutdown finished")

	handler := server.Handler.(*handler)
	handler.close()

	return server.Server.Shutdown(context.Background())
}
