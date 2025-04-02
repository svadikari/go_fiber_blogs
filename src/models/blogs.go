package models

import (
	"time"
)

type Blog struct {
	Id       uint64 `json:"id" gorm:"id primary_key"`
	Title    string `json:"title" form:"title" validate:"required,min=5,max=100" gorm:"title" message:"title is required and should be between 5 and 100 characters length"`
	Content  string `json:"content" form:"content" validate:"required" gorm:"content"`
	AuthorId uint64 `json:"author_id" gorm:"author_id"`
	Author   User   `json:"author" validate:"required=false" gorm:"foreignKey:AuthorId;constraint:OnUpdate:CASCADE,
OnDelete:CASCADE;"`
	CreatedAt time.Time `json:"created_at" gorm:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"updated_at"`
}
