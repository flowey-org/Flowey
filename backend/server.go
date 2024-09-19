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
	var f FloweyServer
	f.Addr = address
	f.Handler = Handler{}
	return &f
}

func (f *FloweyServer) ListenAndServe() error {
	listener, err := net.Listen("tcp", f.Addr)
	if err != nil {
		return err
	}
	log.Printf("listening at %v", f.Addr)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	errs := make(chan error, 1)
	go func() {
		errs <- f.Serve(listener)
	}()

	select {
	case <-sigint:
		fmt.Print("\rinterrupting...\n")
	case err := <-errs:
		log.Printf("failed to serve: %w", err)
	}

	return f.Shutdown(context.Background())

}
