package dto

type (
	StartDTO struct {
		TotalPlayer int `form:"total-player" validate:"required,gte=2,lte=4"`
	}
)
