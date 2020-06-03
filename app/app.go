package app

import (
	"github.com/dqn/tubekids/dbexec"
	"github.com/labstack/echo/v4"
)

type App struct {
	Server *echo.Echo
	DBX    *dbexec.DBExecutor
}

func (a *App) routes() {
	a.Server.GET("/api/channels/:id", a.getChannel)
}

func (a *App) Start(address string) {
	a.routes()
	a.Server.Logger.Fatal(a.Server.Start(address))
}
