package validators

import "github.com/go-playground/validator/v10"

var V = validator.New()

type RegisterDTO struct {
	Username string `json:"username" validate:"required,alphanum,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TaskDTO struct {
	Title       string `json:"title" validate:"required,max=255"`
	Description string `json:"description" validate:"max=5000"`
	Status      string `json:"status" validate:"oneof=pending completed"`
}
