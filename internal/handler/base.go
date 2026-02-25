package handler

import "github.com/mistic0xb/smolurl/internal/server"

type Handler struct {
	server *server.Server
}

func NewHandler(s *server.Server) Handler {
	return Handler{server: s}
}
