package main

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type Handler struct{}

func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	options := websocket.AcceptOptions{InsecureSkipVerify: true}
	connection, err := websocket.Accept(writer, request, &options)
	if err != nil {
		log.Println(err)
		return
	}
	// At this point the client isn't listening anymore
	defer connection.CloseNow()
	for {
		err = handleFrame(connection)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			log.Printf("failed to respond to %v: %v", request.RemoteAddr, err)
			return
		}
	}
}

func handleFrame(connection *websocket.Conn) error {
	ctx := context.TODO()

	messageType, reader, err := connection.Reader(ctx)
	if err != nil {
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
