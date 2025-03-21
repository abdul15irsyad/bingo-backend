package handler

import (
	"bingo/lib"
	"bingo/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(us *service.UserService) *UserHandler {
	lib.Logger.Info("NewUserHandler initialized")
	return &UserHandler{us}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users := h.userService.GetUsers()
	c.JSON(http.StatusOK, gin.H{
		"message": "get users",
		"data":    users,
	})
}

func (h *UserHandler) GetCountUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "get count user",
		"data":    len(h.userService.Users),
	})
}
