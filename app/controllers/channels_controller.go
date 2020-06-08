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
	id := ctx.Param("id")
	channel, err := c.DBX.Channels.Find(id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, channel)
}
