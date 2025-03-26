package handler

import (
	"bingo/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	userService *service.UserService
}

func NewProfileHandler(userService *service.UserService) *ProfileHandler {
	return &ProfileHandler{userService}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	authUser, ok := c.Get("authUser")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile",
		"data":    authUser,
	})
}
