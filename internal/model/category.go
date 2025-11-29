package model

import "time"

type Category struct {
	ID        uint64    `json:"id" db:"id"`
	UserID    uint64    `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	GroupName string    `json:"group_name" db:"group_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateCategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	GroupName string `json:"group_name"`
}

type UpdateCategoryRequest struct {
	Name      *string `json:"name"`
	GroupName *string `json:"group_name"`
}

type UpdateCategoryRecord struct {
	Name      *string
	GroupName *string
	UpdatedAt time.Time
}

type CreateCategoryRecord struct {
	UserID    uint64
	Name      string
	GroupName string
	CreatedAt time.Time
	UpdatedAt time.Time
}
