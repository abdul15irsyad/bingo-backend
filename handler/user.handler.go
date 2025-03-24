package handler

import (
	"bingo/dto"
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
	var getUsersDto dto.GetUsersDto
	c.ShouldBind(&getUsersDto)
	validationErrors := lib.Validate(getUsersDto)
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
			"code":    "VALIDATION_ERROR",
			"errors":  validationErrors,
		})
		return
	}

	users, err := h.userService.GetPaginatedUsers(service.GetPaginatedUsersDto{
		Page:   getUsersDto.Page,
		Limit:  getUsersDto.Limit,
		Search: getUsersDto.Search,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get users",
		"data":    users,
	})
}
