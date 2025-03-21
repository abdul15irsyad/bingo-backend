package routes

import (
	"bingo/handler"
	"bingo/lib"

	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	authHandler *handler.AuthHandler
}

func NewAuthRoute(ah *handler.AuthHandler) *AuthRoute {
	lib.Logger.Info("NewAuthRoute initialized")
	return &AuthRoute{ah}
}

func (r *AuthRoute) InitAuthRoute(router *gin.Engine) {
	auth := router.Group("/auth")
	auth.POST("/register", r.authHandler.Register)
}
