package routes

import (
	"bingo/handler"
	"bingo/lib"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	userHandler *handler.UserHandler
}

func NewUserRoute(uh *handler.UserHandler) *UserRoute {
	lib.Logger.Info("NewUserRoute initialized")
	return &UserRoute{uh}
}

func (r *UserRoute) InitUserRoute(router *gin.Engine) {
	users := router.Group("/users")
	users.GET("/", r.userHandler.GetUsers)
	users.GET("/count", r.userHandler.GetCountUser)
}
