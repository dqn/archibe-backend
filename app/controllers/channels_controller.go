package controllers

import (
	"net/http"

	"github.com/dqn/archibe/dbexec"
	"github.com/labstack/echo/v4"
)

type ChannelsController struct {
	DBX *dbexec.DBExecutor
}

func (c *ChannelsController) GetChannels(ctx echo.Context) error {
	q := ctx.QueryParam("q")
	limit := parseUintWithDefault(ctx.QueryParam("limit"), 30)
	offset := parseUintWithDefault(ctx.QueryParam("offset"), 0)

	channels, err := c.DBX.Channels.FindByQuery(&dbexec.ChannelsQuery{
		Q:      q,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, channels)
}

func (c *ChannelsController) GetChannel(ctx echo.Context) error {
	id := ctx.Param("id")
	channel, err := c.DBX.Channels.Find(id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, channel)
}
