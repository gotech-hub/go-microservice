package bootstrap

import "go-source/api/http/handlers"

type Handlers struct {
	Handler *handlers.Handler
}

func NewHandlers(services *Services) *Handlers {
	return &Handlers{
		Handler: handlers.NewHandler(),
	}
}
