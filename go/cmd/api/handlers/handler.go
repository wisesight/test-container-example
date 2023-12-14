package handlers

import "github.com/wisesight/test-container-example/internal/usecases"

type EchoHandler struct {
	usecases usecases.Usecases
}

func NewEchoHandler(usecases usecases.Usecases) *EchoHandler {
	return &EchoHandler{
		usecases: usecases,
	}
}
