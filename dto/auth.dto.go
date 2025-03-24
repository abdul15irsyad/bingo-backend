package dto

type RegisterDTO struct {
	Name            string `json:"name" validate:"required"`
	Username        string `json:"username" validate:"required,username"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}
