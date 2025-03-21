package routes

import (
	"bingo/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RootRoute struct{}

func NewRootRoute() *RootRoute {
	lib.Logger.Info("NewRootRoute initialized")
	return &RootRoute{}
}

func (r *RootRoute) InitRootRoute(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "bingo backend",
		})
	})
}
