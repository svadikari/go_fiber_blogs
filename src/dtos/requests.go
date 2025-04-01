package dtos

type LoginRequest struct {
	UserName string `json:"username" validate:"required,min=5,max=10"`
	Password string `json:"password" validate:"required,min=5,max=10"`
}
