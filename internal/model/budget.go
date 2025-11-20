package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type Budget struct {
	ID         uint64          `json:"id" db:"id"`
	UserID     uint64          `json:"user_id" db:"user_id"`
	CategoryID uint64          `json:"category_id" db:"category_id"`
	Year       uint            `json:"year" db:"year"`
	Month      uint            `json:"month" db:"month"`
	Budgeted   decimal.Decimal `json:"budgeted" db:"budgeted"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}

type BudgetAllocationRequest struct {
	CategoryID uint64          `json:"category_id" binding:"required"`
	Year       uint            `json:"year" binding:"required"`
	Month      uint            `json:"month" binding:"required"`
	Amount     decimal.Decimal `json:"amount" binding:"required"`
}

type BillPaymentRequest struct {
	BillInstanceID uint64          `json:"bill_instance_id" binding:"required"`
	Amount         decimal.Decimal `json:"amount" binding:"required"`
	TransactionID  *uint64         `json:"transaction_id,omitempty"`
}
