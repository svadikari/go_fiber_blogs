package dtos

type LoginRequest struct {
	UserName string `json:"username" validate:"required,min=5,max=10"`
	Password string `json:"password" validate:"required,min=5,max=10"`
}

type BlogRequest struct {
	Title   string `json:"title" form:"title" validate:"required,min=5,max=100" gorm:"title" message:"title is required and should be between 5 and 100 characters length"`
	Content string `json:"content" form:"content" validate:"required" gorm:"content"`
}

type UserProfile struct {
	FistName string `json:"firstName" form:"firstName" validate:"required,min=5,max=100" message:"firstName is required and should be between 5 and 100 characters length"`
	LastName string `json:"lastName" form:"lastName" validate:"required,min=5,max=100" message:"lastName is required and should be between 5 and 100 characters length"`
	Phone    string `json:"phone" form:"phone" validate:"required,numeric,len=10"`
}

type UserRequest struct {
	UserProfile
	LoginRequest
}
