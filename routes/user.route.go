package routes

import (
	"bingo/handler"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	users.GET("/count", handler.GetCountUser)
}
