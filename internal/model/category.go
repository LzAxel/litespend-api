package model

import "time"

type CategoryType uint8

const (
	CategoryTypeIncome  CategoryType = 1
	CategoryTypeExpense CategoryType = 2
)

type TransactionCategory struct {
	ID        uint64       `json:"id" db:"id"`
	UserID    uint64       `json:"user_id" db:"user_id"`
	Name      string       `json:"name" db:"name"`
	Type      CategoryType `json:"type" db:"type"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
}

type CreateCategoryRequest struct {
	Name string       `json:"name" binding:"required"`
	Type CategoryType `json:"type"`
}

type UpdateCategoryRequest struct {
	Name *string       `json:"name"`
	Type *CategoryType `json:"type"`
}

type CreateCategoryRecord struct {
	UserID    uint64
	Name      string
	Type      CategoryType
	CreatedAt time.Time
}
