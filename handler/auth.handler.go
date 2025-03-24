package handler

import (
	"bingo/config"
	"bingo/dto"
	"bingo/lib"
	"bingo/service"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(us *service.UserService) *AuthHandler {
	lib.Logger.Info("NewAuthHandler initialized")
	return &AuthHandler{us}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var registerDTO dto.RegisterDTO
	c.ShouldBindJSON(&registerDTO)
	validationErrors := lib.Validate(registerDTO)
	if len(validationErrors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
			"errors":  validationErrors,
		})
		return
	}

	trimmedUsername := strings.TrimSpace(registerDTO.Username)
	_, err := h.userService.GetUserByUsername(trimmedUsername)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	} else if err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
			"code":    "VALIDATION_ERROR",
			"errors": []map[string]any{
				{
					"field":   "username",
					"message": "Username already exist",
				},
			},
		})
	}

	newUser, err := h.userService.CreateUser(service.CreateUserDTO{
		Name:     registerDTO.Name,
		Username: trimmedUsername,
		Email:    registerDTO.Email,
		Password: registerDTO.Password,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	// create access token
	accessToken, err := lib.CreateJWTToken(&jwt.MapClaims{
		"sub": newUser.Id,
	})
	if err != nil {
		fmt.Println(err)
	}

	c.SetCookie("accessToken", accessToken, 0, "/", config.CookieDomain, false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Bad Request",
		"data":    newUser,
	})
}

func (h *AuthHandler) Profile(c *gin.Context) {

}
