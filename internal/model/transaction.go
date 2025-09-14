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
	ID          int             `json:"id" db:"id"`
	UserID      int             `json:"user_id" db:"user_id"`
	CategoryID  int             `json:"category_id" db:"category_id"`
	GoalID      *int            `json:"goal_id,omitempty" db:"goal_id"`
	Description string          `json:"description" db:"description"`
	Amount      decimal.Decimal `json:"amount" db:"amount"`
	Type        TransactionType `json:"type" db:"type"`
	DateTime    time.Time       `json:"date_time" db:"date_time"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
}

type CreateTransactionRecord struct {
	UserID      int
	CategoryID  int
	GoalID      *int
	Description string
	Amount      decimal.Decimal
	Type        TransactionType
	DateTime    time.Time
	CreatedAt   time.Time
}

type UpdateTransactionRecord struct {
	CategoryID  int
	GoalID      *int
	Description string
	Amount      decimal.Decimal
	Type        TransactionType
	DateTime    time.Time
}
