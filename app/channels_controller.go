package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *App) getChannel(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
