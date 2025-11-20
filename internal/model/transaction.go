package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	ID             uint64          `json:"id" db:"id"`
	UserID         uint64          `json:"user_id" db:"user_id"`
	CategoryID     *uint64         `json:"category_id,omitempty" db:"category_id"`
	BillInstanceID *uint64         `json:"bill_instance_id,omitempty" db:"bill_instance_id"`
	Amount         decimal.Decimal `json:"amount" db:"amount"`
	Date           time.Time       `json:"date" db:"date"`
	Description    string          `json:"description" db:"description"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
}

type CreateTransactionRequest struct {
	CategoryID     *uint64         `json:"category_id,omitempty"`
	Amount         decimal.Decimal `json:"amount"`
	Date           time.Time       `json:"date"`
	Description    string          `json:"description"`
	BillInstanceID *uint64         `json:"bill_instance_id,omitempty"`
}

type UpdateTransactionRequest struct {
	BillInstanceID *uint64          `json:"bill_instance_id,omitempty"`
	CategoryID     *uint64          `json:"category_id,omitempty"`
	Amount         *decimal.Decimal `json:"amount,omitempty"`
	Date           *time.Time       `json:"date,omitempty"`
	Description    *string          `json:"description,omitempty"`
}

type PaginatedTransactionsResponse = PaginatedResponse[Transaction]
