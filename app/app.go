package app

import (
	"github.com/dqn/tubekids/app/controllers"
	"github.com/dqn/tubekids/dbexec"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	Server *echo.Echo
	DBX    *dbexec.DBExecutor
}

func (a *App) routes() {
	channels := controllers.ChannelsController{DBX: a.DBX}
	a.Server.GET("/api/channels/:id", channels.GetChannel)

	channelChats := controllers.ChatsController{DBX: a.DBX}
	a.Server.GET("/api/chats", channelChats.GetChats)
}

func (a *App) Start(address string) {
	a.Server.Use(middleware.Logger())
	a.Server.Use(middleware.Recover())

	a.routes()

	a.Server.Logger.Fatal(a.Server.Start(address))
}
