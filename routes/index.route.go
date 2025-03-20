package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "bingo backend",
		})
	})
	AuthRoutes(r)
	UserRoutes(r)
}
