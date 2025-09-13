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
	ID          uint64          `json:"id" db:"id"`
	CategoryID  uint64          `json:"category_id" db:"category_id"`
	GoalID      *uint64         `json:"goal_id,omitempty" db:"goal_id"`
	Description string          `json:"description" db:"description"`
	Amount      decimal.Decimal `json:"amount" db:"amount"`
	Type        TransactionType `json:"type" db:"type"`
	DateTime    time.Time       `json:"date_time" db:"date_time"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
}
