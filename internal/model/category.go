package model

import "time"

type TransactionCategory struct {
	ID        uint64          `json:"id" db:"id"`
	Name      string          `json:"name" db:"name"`
	Type      TransactionType `json:"type" db:"type"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}
