package middleware

import (
	"bingo/lib"

	"github.com/gin-gonic/gin"
)

type CorsMiddleware struct{}

func NewCorsMiddleware() *CorsMiddleware {
	lib.Logger.Info("CorsMiddleware initialized")
	return &CorsMiddleware{}
}

func (m *CorsMiddleware) Handler(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	c.Next()
}
