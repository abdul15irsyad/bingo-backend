package handler

import (
	"bingo/config"
	"bingo/dto"
	"bingo/lib"
	"bingo/model"
	"bingo/service"
	"bingo/util"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	validationErrors := lib.Validate(c, registerDTO)
	if len(validationErrors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
			"errors":  validationErrors,
		})
		return
	}

	if len(h.userService.Users) >= h.userService.MaxUser {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
			"code":    "MAX_USER_EXCEED",
			"error":   "Max User Exceed",
		})
		return
	}

	trimmedUsername := strings.TrimSpace(registerDTO.Username)
	if userWithUsernameExist := util.FindSlice(&h.userService.Users, func(user model.User) bool {
		return util.Slugify(user.Username) == util.Slugify(trimmedUsername)
	}); userWithUsernameExist != nil {
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
		return
	}

	newUser := h.userService.AddUser(service.UserDTO{
		Username: trimmedUsername,
	})

	fmt.Println(
		strings.Join(
			util.MapSlice(
				h.userService.Users,
				func(user model.User) string {
					return user.Username
				}),
			", ",
		),
	)

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
