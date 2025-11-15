package model

import "time"

type TransactionCategory struct {
	ID        uint64          `json:"id" db:"id"`
	UserID    uint64          `json:"user_id" db:"user_id"`
	Name      string          `json:"name" db:"name"`
	Type      TransactionType `json:"type" db:"type"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

type CreateCategoryRequest struct {
	Name string          `json:"name" binding:"required"`
	Type TransactionType `json:"type"`
}

type UpdateCategoryRequest struct {
	Name *string          `json:"name"`
	Type *TransactionType `json:"type"`
}

type CreateCategoryRecord struct {
	UserID    uint64
	Name      string
	Type      TransactionType
	CreatedAt time.Time
}
