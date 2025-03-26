package middleware

import (
	"bingo/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorMiddleware struct{}

func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

func (m *ErrorMiddleware) Handler(c *gin.Context) {
	c.Next()
	defer func() {
		if r := recover(); r != nil {
			lib.Logger.Error(r.(error).Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
		}
	}()

	if len(c.Errors) > 0 {
		err := c.Errors.Last().Err
		if err != nil {
			lib.Logger.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
		}
	}
}
