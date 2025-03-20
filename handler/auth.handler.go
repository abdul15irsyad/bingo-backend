package handler

import (
	"bingo/data"
	"bingo/dto"
	"bingo/service"
	"bingo/util"
	"bingo/validation"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var registerDTO dto.RegisterDTO
	c.ShouldBindJSON(&registerDTO)
	validationErrors := validation.Validate(c, registerDTO)
	if len(validationErrors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
			"errors":  validationErrors,
		})
		return
	}

	if len(data.Users) >= data.MaxUser {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
			"code":    "MAX_USER_EXCEED",
			"error":   "Max User Exceed",
		})
		return
	}

	trimmedUsername := strings.TrimSpace(registerDTO.Username)
	if userWithUsernameExist := util.FindSlice(&data.Users, func(user data.User) bool {
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

	newUser := service.AddUser(service.AddUserDTO{
		Username: trimmedUsername,
	})

	// fmt.Println(strings.Join(util.MapSlice(data.Users, func(user data.User) string { return user.Username }), ", "))

	c.JSON(http.StatusOK, gin.H{
		"message": "Bad Request",
		"data":    newUser,
	})
}
