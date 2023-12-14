package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (eh *EchoHandler) GetPostByID(echoCtx echo.Context) error {
	id := echoCtx.Param("id")
	post, err := eh.usecases.GetPostByID(id)
	if err != nil {
		return err
	}
	return echoCtx.JSON(http.StatusOK, post)
}
