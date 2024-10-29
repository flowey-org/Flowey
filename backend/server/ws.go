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

type connections struct {
	data  map[*connection]bool
	mutex sync.Mutex
}

func newConnections() connections {
	return connections{data: make(map[*connection]bool)}
}

func (connections *connections) store(key *connection, value bool) {
	connections.mutex.Lock()
	connections.data[key] = value
	connections.mutex.Unlock()
}

func (connections *connections) delete(key *connection) {
	connections.mutex.Lock()
	delete(connections.data, key)
	connections.mutex.Unlock()
}

func (connections *connections) close() {
	connections.mutex.Lock()
	defer connections.mutex.Unlock()

	for connection := range connections.data {
		connection.Close(websocket.StatusNormalClosure, "server shutting down")
	}
}

type wsHandler struct {
	connections connections
	waitGroup   sync.WaitGroup
}

func newWsHandler() *wsHandler {
	return &wsHandler{
		connections: newConnections(),
	}
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

	handler.connections.store(&connection, true)
	defer handler.connections.delete(&connection)

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

	_, err = db.AuthenticateBySessionKey(sessionKey.Value)
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
	handler.connections.close()
	handler.waitGroup.Wait()
}
