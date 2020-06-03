package app

import (
	"github.com/dqn/tubekids/app/controllers"
	"github.com/dqn/tubekids/dbexec"
	"github.com/labstack/echo/v4"
)

type App struct {
	Server *echo.Echo
	DBX    *dbexec.DBExecutor
}

func (a *App) routes() {
	channels := controllers.ChannelsController{DBX: a.DBX}
	a.Server.GET("/api/channels/:id", channels.GetChannel)
}

func (a *App) Start(address string) {
	a.routes()
	a.Server.Logger.Fatal(a.Server.Start(address))
}
