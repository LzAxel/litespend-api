package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type TransactionType uint8

const (
	TransactionTypeIncome TransactionType = iota
	TransactionTypeExpanse
)

type Transaction struct {
	ID                  uint64          `json:"id" db:"id"`
	UserID              uint64          `json:"user_id" db:"user_id"`
	CategoryID          uint64          `json:"category_id" db:"category_id"`
	GoalID              *uint64         `json:"goal_id,omitempty" db:"goal_id"`
	PrescribedExpanseID *uint64         `json:"prescribed_expanse_id,omitempty" db:"prescribed_expanse_id"`
	Description         string          `json:"description" db:"description"`
	Amount              decimal.Decimal `json:"amount" db:"amount"`
	Type                TransactionType `json:"type" db:"type"`
	DateTime            time.Time       `json:"date_time" db:"date_time"`
	CreatedAt           time.Time       `json:"created_at" db:"created_at"`
}

type CreateTransactionRequest struct {
	CategoryID          uint64          `json:"category_id"`
	GoalID              *uint64         `json:"goal_id,omitempty"`
	PrescribedExpanseID *uint64         `json:"prescribed_expanse_id,omitempty"`
	Description         string          `json:"description"`
	Amount              decimal.Decimal `json:"amount"`
	Type                TransactionType `json:"type"`
	DateTime            time.Time       `json:"date_time"`
}

type UpdateTransactionRequest struct {
	CategoryID          *uint64          `json:"category_id"`
	GoalID              *uint64          `json:"goal_id,omitempty"`
	PrescribedExpanseID *uint64          `json:"prescribed_expanse_id,omitempty"`
	Description         *string          `json:"description"`
	Amount              *decimal.Decimal `json:"amount"`
	Type                *TransactionType `json:"type"`
	DateTime            *time.Time       `json:"date_time"`
}

type PaginatedTransactionsResponse = PaginatedResponse[Transaction]
