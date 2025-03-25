package routes

import (
	"bingo/handler"
	"bingo/lib"

	"github.com/gin-gonic/gin"
)

type SocketRoute struct {
	socketHandler *handler.SocketHandler
}

func NewSocketRoute(socketHandler *handler.SocketHandler) *SocketRoute {
	lib.Logger.Info("NewSocketRoute initialized")
	return &SocketRoute{socketHandler}
}

func (r *SocketRoute) InitSocketRoute(router *gin.Engine) {
	socket := router.Group("/socket")
	socket.GET("/start-game", r.socketHandler.StartGameHandler)
}
