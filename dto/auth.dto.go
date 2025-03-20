package dto

type RegisterDTO struct {
	Username string `json:"username" validate:"required,username"`
}
