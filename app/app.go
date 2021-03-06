package app

import (
	"github.com/dqn/archibe/app/controllers"
	"github.com/dqn/archibe/dbexec"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	Server *echo.Echo
	DBX    *dbexec.DBExecutor
}

func (a *App) routes() {
	channels := controllers.ChannelsController{DBX: a.DBX}
	a.Server.GET("/api/channels", channels.GetChannels)
	a.Server.GET("/api/channels/:id", channels.GetChannel)

	chats := controllers.ChatsController{DBX: a.DBX}
	a.Server.GET("/api/chats", chats.GetChats)

	videos := controllers.VideosController{DBX: a.DBX}
	a.Server.GET("/api/videos", videos.GetVideos)
	a.Server.GET("/api/videos/:id", videos.GetVideo)
}

func (a *App) Start(address string) {
	a.Server.Use(middleware.Logger())
	a.Server.Use(middleware.Recover())

	a.routes()

	a.Server.Logger.Fatal(a.Server.Start(address))
}
