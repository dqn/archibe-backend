package controllers

import (
	"net/http"

	"github.com/dqn/tubekids/dbexec"
	"github.com/labstack/echo/v4"
)

type ChatsController struct {
	DBX *dbexec.DBExecutor
}

func (c *ChatsController) GetChats(ctx echo.Context) error {
	channelID := ctx.QueryParam("channel_id")
	videoID := ctx.QueryParam("video_id")
	order := ctx.QueryParam("order")
	q := ctx.QueryParam("q")
	limit := parseUintWithDefault(ctx.QueryParam("limit"), 30)
	offset := parseUintWithDefault(ctx.QueryParam("offset"), 0)

	chats, err := c.DBX.Chats.FindByQuery(&dbexec.ChatsQuery{
		Q:       q,
		Channel: channelID,
		Video:   videoID,
		Order:   order,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, chats)
}
