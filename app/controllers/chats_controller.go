package controllers

import (
	"net/http"
	"strconv"

	"github.com/dqn/tubekids/dbexec"
	"github.com/labstack/echo/v4"
)

type ChatsController struct {
	DBX *dbexec.DBExecutor
}

func parseUintWithDefault(str string, dflt uint64) uint64 {
	if v, err := strconv.ParseUint(str, 10, 64); err == nil {
		return v
	} else {
		return dflt
	}
}

func (c *ChatsController) GetChats(ctx echo.Context) error {
	channelID := ctx.QueryParam("channel_id")
	limit := parseUintWithDefault(ctx.QueryParam("limit"), 30)
	offset := parseUintWithDefault(ctx.QueryParam("offset"), 0)

	chats, err := c.DBX.Chats.FindByQuery(&dbexec.ChatsQuery{
		Channel: channelID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, chats)
}
