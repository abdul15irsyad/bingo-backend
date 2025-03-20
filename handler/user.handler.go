package handler

import (
	"bingo/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCountUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "get count user",
		"data":    len(data.Users),
	})
}
