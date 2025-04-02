package models

import (
	"time"
)

type User struct {
	Id        uint64    `json:"id" gorm:"id primary_key"`
	FistName  string    `json:"firstName" validate:"required,min=5,max=100" gorm:"first_name" message:"firstName is required and should be between 5 and 100 characters length"`
	LastName  string    `json:"lastName" validate:"required,min=5,max=100" gorm:"last_name" message:"lastName is required and should be between 5 and 100 characters length"`
	Phone     string    `json:"phone" validate:"required" gorm:"phone"`
	UserName  string    `json:"userName" validate:"required" gorm:"username"`
	Password  string    `json:"-" validate:"required" gorm:"password"`
	CreatedAt time.Time `json:"created_at" gorm:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"updated_at"`
	Blogs     []Blog    `json:"blogs" gorm:"foreignKey:AuthorId"`
}
