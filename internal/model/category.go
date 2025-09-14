package model

import "time"

type TransactionCategory struct {
	ID        int             `json:"id" db:"id"`
	UserID    int             `json:"user_id" db:"user_id"`
	Name      string          `json:"name" db:"name"`
	Type      TransactionType `json:"type" db:"type"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

type CreateTransactionCategoryRecord struct {
	Name      string
	UserID    int
	Type      TransactionType
	CreatedAt time.Time
}

type UpdateTransactionCategoryRecord struct {
	Name string
	Type TransactionType
}
