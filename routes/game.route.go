package routes

import (
	"bingo/handler"
	"bingo/lib"

	"github.com/gin-gonic/gin"
)

type GameRoute struct {
	gameHandler *handler.GameHandler
}

func NewGameRoute(gameHandler *handler.GameHandler) *GameRoute {
	lib.Logger.Info("NewGameRoute initialized")
	return &GameRoute{gameHandler}
}

func (r *GameRoute) InitGameRoute(router *gin.Engine) {
	game := router.Group("/game")
	game.GET("/start", r.gameHandler.Start)
}
