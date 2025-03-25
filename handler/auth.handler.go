package handler

import (
	"bingo/config"
	"bingo/dto"
	"bingo/lib"
	"bingo/service"
	"bingo/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	lib.Logger.Info("NewAuthHandler initialized")
	return &AuthHandler{userService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginDTO dto.LoginDTO
	c.ShouldBindJSON(&loginDTO)
	validationErrors := lib.Validate(loginDTO)
	if len(validationErrors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
			"errors":  validationErrors,
		})
		return
	}

	authUser, err := h.userService.GetUserByUsernameOrEmail(loginDTO.UsernameOrEmail)
	if err != nil {
		util.ComparePassword("some password", loginDTO.Password)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "email or password is incorrect",
		})
		return
	}

	correctPassword, err := util.ComparePassword(authUser.Password, loginDTO.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	if !correctPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email or password is incorrect",
		})
		return
	}

	// create jwt
	accessToken, err := lib.CreateJWT(&jwt.MapClaims{
		"sub": authUser.Id,
	})
	if err != nil {
		return
	}

	c.SetCookie("access_token", accessToken, 0, "/", config.CookieDomain, false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login",
	})
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Register",
		"data":    newUser,
	})
}

func (h *AuthHandler) Profile(c *gin.Context) {

}
