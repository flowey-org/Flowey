package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/coder/websocket"
)

type connection struct {
	*websocket.Conn
	writer  http.ResponseWriter
	request *http.Request
}

func (connection *connection) handleFrame(ctx context.Context) error {
	messageType, message, err := connection.Read(ctx)
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

	_ = messageType
	_ = message

	return nil
}

type handler struct {
	connections sync.Map
	waitGroup sync.WaitGroup
}

func (handler *handler) handle(writer http.ResponseWriter, request *http.Request) error {
	defer handler.waitGroup.Done()

	options := websocket.AcceptOptions{InsecureSkipVerify: true}
	conn, err := websocket.Accept(writer, request, &options)
	if err != nil {
		return err
	}

	connection := connection{conn, writer, request}
	defer connection.CloseNow()

	handler.connections.Store(connection, true)
	defer handler.connections.Delete(connection)

	log.Printf("opened a connection with %v", request.RemoteAddr)

	ctx := context.Background()
	for {
		err := connection.handleFrame(ctx)
		if err != nil {
			return err
		}
	}
}

func (handler *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler.waitGroup.Add(1)

	err := handler.handle(writer, request)
	if err != nil {
		log.Println(err)
		return
	}
}

func (handler *handler) close() {
	for key, _ := range handler.connections.Range {
		connection, ok := key.(connection)
		if !ok {
			return
		}

		connection.Close(websocket.StatusNormalClosure, "server shutting down")
	}

	handler.waitGroup.Wait()
}
