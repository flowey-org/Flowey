package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type connection struct {
	*websocket.Conn
	writer  http.ResponseWriter
	request *http.Request
}

func handleConnection(writer http.ResponseWriter, request *http.Request) error {
	options := websocket.AcceptOptions{InsecureSkipVerify: true}
	conn, err := websocket.Accept(writer, request, &options)
	if err != nil {
		return err
	}

	connection := connection{conn, writer, request}
	defer connection.CloseNow()

	log.Printf("opened a connection with %v", request.RemoteAddr)

	for {
		err := connection.handleFrame()
		if err != nil {
			return err
		}
	}
}

func (connection *connection) handleFrame() error {
	ctx := context.Background()

	messageType, reader, err := connection.Reader(ctx)
	if err != nil {
		var closeError websocket.CloseError
		if errors.As(err, &closeError) {
			return fmt.Errorf(
				"received a close frame from %v (%d, %q)",
				connection.request.RemoteAddr, closeError.Code, closeError.Reason,
			)
		}

		return err
	}

	writer, err := connection.Writer(ctx, messageType)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		return err
	}

	return writer.Close()
}

type Handler struct{}

func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	err := handleConnection(writer, request)
	if err != nil {
		log.Println(err)
		return
	}
}
