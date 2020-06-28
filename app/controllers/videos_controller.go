package controllers

import (
	"net/http"

	"github.com/dqn/tubekids/dbexec"
	"github.com/labstack/echo/v4"
)

type VideosController struct {
	DBX *dbexec.DBExecutor
}

func (c *VideosController) GetVideo(ctx echo.Context) error {
	id := ctx.Param("id")

	video, err := c.DBX.Videos.Find(id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, video)
}

func (c *VideosController) GetVideos(ctx echo.Context) error {
	q := ctx.QueryParam("q")
	channelID := ctx.QueryParam("channel_id")
	limit := parseUintWithDefault(ctx.QueryParam("limit"), 30)
	offset := parseUintWithDefault(ctx.QueryParam("offset"), 0)
	order := ctx.QueryParam("order")

	channels, err := c.DBX.Videos.FindByQuery(&dbexec.VideosQuery{
		Q:       q,
		Channel: channelID,
		Limit:   limit,
		Offset:  offset,
		Order:   order,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, channels)
}
