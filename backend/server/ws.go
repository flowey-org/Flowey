package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/coder/websocket"

	"flowey/db"
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

type wsHandler struct {
	connections sync.Map
	waitGroup   sync.WaitGroup
}

func (handler *wsHandler) handle(writer http.ResponseWriter, request *http.Request) error {
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

func (handler *wsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	sessionKey, err := request.Cookie(sessionKeyCookieName)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = db.Authorize(sessionKey.Value)
	if err != nil {
		switch err {
		case db.Unathorized:
			writer.WriteHeader(http.StatusUnauthorized)
		case db.InternalServerError:
			writer.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	handler.waitGroup.Add(1)

	err = handler.handle(writer, request)
	if err != nil {
		log.Println(err)
		return
	}
}

func (handler *wsHandler) close() {
	for key, _ := range handler.connections.Range {
		connection, ok := key.(connection)
		if !ok {
			return
		}

		connection.Close(websocket.StatusNormalClosure, "server shutting down")
	}

	handler.waitGroup.Wait()
}
