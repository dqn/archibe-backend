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
