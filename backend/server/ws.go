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

type userConnections = map[*connection]bool

func newUserConnections() userConnections {
	return make(map[*connection]bool)
}

type connections struct {
	dict  map[db.UserID]userConnections
	mutex sync.RWMutex
}

func newConnections() connections {
	return connections{dict: make(map[db.UserID]userConnections)}
}

func (connections *connections) store(userID db.UserID, connection *connection) {
	connections.mutex.Lock()
	defer connections.mutex.Unlock()

	if userConnections, ok := connections.dict[userID]; ok {
		userConnections[connection] = true
		return
	}

	userConnections := newUserConnections()
	userConnections[connection] = true
	connections.dict[userID] = userConnections
}

func (connections *connections) delete(userID db.UserID, connection *connection) {
	connections.mutex.Lock()
	defer connections.mutex.Unlock()

	if userConnections, ok := connections.dict[userID]; ok {
		delete(userConnections, connection)
	}
}

func (connections *connections) close() {
	connections.mutex.RLock()
	defer connections.mutex.RUnlock()

	for _, userConnections := range connections.dict {
		for connection := range userConnections {
			connection.Close(websocket.StatusNormalClosure, "server shutting down")
		}
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

func (handler *wsHandler) handle(userID db.UserID, writer http.ResponseWriter, request *http.Request) error {
	defer handler.waitGroup.Done()

	options := websocket.AcceptOptions{InsecureSkipVerify: true}
	conn, err := websocket.Accept(writer, request, &options)
	if err != nil {
		return err
	}

	connection := connection{conn, writer, request}
	defer connection.CloseNow()

	handler.connections.store(userID, &connection)
	defer handler.connections.delete(userID, &connection)

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

	userID, err := db.AuthenticateBySessionKey(sessionKey.Value)
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

	err = handler.handle(userID, writer, request)
	if err != nil {
		log.Println(err)
		return
	}
}

func (handler *wsHandler) close() {
	handler.connections.close()
	handler.waitGroup.Wait()
}
