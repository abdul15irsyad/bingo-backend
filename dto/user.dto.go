package dto

type (
	GetUserDto struct {
		Id string `validate:"required,uuid"`
	}

	GetUsersDto struct {
		Page   int     `json:"page" validate:"required,number,gte=1"`
		Limit  int     `json:"limit" validate:"required,number"`
		Search *string `json:"search" validate:"omitempty"`
	}
)
