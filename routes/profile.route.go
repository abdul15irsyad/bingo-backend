package routes

import (
	"bingo/handler"
	"bingo/lib"

	"github.com/gin-gonic/gin"
)

type ProfileRoute struct {
	profileHandler *handler.ProfileHandler
}

func NewProfileRoute(ah *handler.ProfileHandler) *ProfileRoute {
	lib.Logger.Info("NewProfileRoute initialized")
	return &ProfileRoute{ah}
}

func (r *ProfileRoute) InitProfileRoute(router *gin.Engine) {
	profile := router.Group("/profile")
	profile.GET("/", r.profileHandler.GetProfile)
}
