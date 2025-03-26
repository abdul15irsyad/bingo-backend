package middleware

import (
	"bingo/lib"
	"bingo/service"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JWTMiddleware struct {
	userService *service.UserService
}

func NewJWTMiddleware(us *service.UserService) *JWTMiddleware {
	return &JWTMiddleware{us}
}

func (m *JWTMiddleware) Handler(c *gin.Context) {
	// accessToken, err := c.Cookie("access_token")
	// if err != nil || accessToken == "" {
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
	// 		"message": "Invalid Credential",
	// 	})
	// 	return
	// }
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid Credential",
		})
		return
	}
	accessToken := strings.Split(authorization, " ")[1]

	sub, err := lib.ParseJWT(accessToken)
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
				"message": "Invalid Credential",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid Credential",
		})
		return
	}

	// check to database
	userId, err := uuid.Parse(sub)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid Credential",
		})
		return
	}
	authUser, err := m.userService.GetUser(userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid Credential",
			})
			return
		}
		c.Error(err)
		return
	}

	c.Set("authUser", authUser)
	c.Next()
}
