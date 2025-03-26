package handler

import (
	"bingo/dto"
	"bingo/lib"
	"bingo/service"
	"bingo/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	lib.Logger.Info("NewAuthHandler initialized")
	return &AuthHandler{userService}
}

func (h *AuthHandler) GuestLogin(c *gin.Context) {

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
		// util.ComparePassword("some password", loginDTO.Password)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "email or password is incorrect",
		})
		return
	}

	correctPassword, err := util.ComparePassword(authUser.Password, loginDTO.Password)
	if err != nil {
		c.Error(err)
		return
	}
	if !correctPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email or password is incorrect",
		})
		return
	}

	// create jwt
	accessToken, err := lib.CreateJWT(authUser.Id.String())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login",
		"data": gin.H{
			"access_token": accessToken,
		},
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
		c.Error(err)
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
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Register",
		"data":    newUser,
	})
}

func (h *AuthHandler) Profile(c *gin.Context) {

}
