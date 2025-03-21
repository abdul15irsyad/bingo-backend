package middleware

import (
	"bingo/lib"
	"bingo/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTMiddleware struct {
	userService *service.UserService
}

func NewJWTMiddleware(us *service.UserService) *JWTMiddleware {
	return &JWTMiddleware{us}
}

func (m *JWTMiddleware) Handler(c *gin.Context) {
	accessToken, err := c.Cookie("accessToken")
	if err != nil || accessToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid credential",
		})
		return
	}

	claims, err := lib.ParseJWTToken(accessToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token expired",
				"code":    "TOKEN_EXPIRED",
			})
			return
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid credential",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	// check to database
	id, _ := claims.GetSubject()
	userId, _ := uuid.Parse(id)
	authUser := m.userService.GetUser(userId)
	if authUser == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid credential",
		})
		return
	}

	c.Set("authUser", authUser)
	c.Next()
}
