package controllers

import (
	"net/http"

	"github.com/dqn/tubekids/dbexec"
	"github.com/labstack/echo/v4"
)

type ChannelsController struct {
	DBX *dbexec.DBExecutor
}

func (c *ChannelsController) GetChannel(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello, World!")
}
