package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	ID         uint64          `json:"id" db:"id"`
	UserID     uint64          `json:"user_id" db:"user_id"`
	CategoryID uint64          `json:"category_id,omitempty" db:"category_id"`
	AccountID  uint64          `json:"account_id" db:"account_id"`
	Note       string          `json:"note" db:"note"`
	Amount     decimal.Decimal `json:"amount" db:"amount"`
	Date       time.Time       `json:"date" db:"date"`
	IsCleared  bool            `json:"is_cleared" db:"cleared"`
	IsApproved bool            `json:"is_approved" db:"approved"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at" db:"updated_at"`
}

type CreateTransactionRequest struct {
	AccountID  uint64          `json:"account_id"`
	CategoryID *uint64         `json:"category_id,omitempty"`
	Amount     decimal.Decimal `json:"amount"`
	Note       string          `json:"note"`
	Date       time.Time       `json:"date"`
	IsCleared  bool            `json:"is_cleared"`
	IsApproved bool            `json:"is_approved"`
}

type CreateTransactionRecord struct {
	UserID     uint64
	AccountID  uint64
	CategoryID *uint64
	Amount     decimal.Decimal
	Note       string
	Date       time.Time
	IsCleared  bool
	IsApproved bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UpdateTransactionRequest struct {
	AccountID  *uint64          `json:"account_id,omitempty"`
	CategoryID *uint64          `json:"category_id,omitempty"`
	Amount     *decimal.Decimal `json:"amount,omitempty"`
	Date       *time.Time       `json:"date,omitempty"`
	Note       *string          `json:"note,omitempty"`
	IsCleared  *bool            `json:"is_cleared,omitempty"`
	IsApproved *bool            `json:"is_approved,omitempty"`
}

type UpdateTransactionRecord struct {
	AccountID  *uint64
	CategoryID *uint64
	Amount     *decimal.Decimal
	Date       *time.Time
	Note       *string
	IsCleared  *bool
	IsApproved *bool
	UpdatedAt  time.Time
}
type PaginatedTransactionsResponse = PaginatedResponse[Transaction]
