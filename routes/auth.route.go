package routes

import (
	"bingo/handler"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	auth.POST("/register", handler.Register)
}
