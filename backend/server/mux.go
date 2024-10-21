package server

import "net/http"

type ServeMux struct {
	http.ServeMux

	ws wsHandler
}

func NewServeMux() *ServeMux {
	mux := ServeMux{}
	mux.Handle("GET /ws", &mux.ws)
	return &mux
}

func (s *ServeMux) close() {
	s.ws.close()
}
