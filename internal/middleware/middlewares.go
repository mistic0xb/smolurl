package middleware

import "github.com/mistic0xb/smolurl/internal/server"

type Middlewares struct {
	Global          *GlobalMiddlewares
	ContextEnhancer *ContextEnhancer
}

func NewMiddleWares(s *server.Server) *Middlewares {
	return &Middlewares{
		Global:          NewGlobalMiddlewares(s),
		ContextEnhancer: NewContextEnhancer(s),
	}
}
