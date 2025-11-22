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

type BudgetDetailed struct {
	ID         uint64          `json:"id" db:"id"`
	UserID     uint64          `json:"user_id" db:"user_id"`
	CategoryID uint64          `json:"category_id" db:"category_id"`
	Year       uint            `json:"year" db:"year"`
	Month      uint            `json:"month" db:"month"`
	Budgeted   decimal.Decimal `json:"budgeted" db:"budgeted"`
	Spent      decimal.Decimal `json:"spent" db:"spent"`
	Remaining  decimal.Decimal `json:"remaining" db:"remaining"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}

type CreateBudgetRecord struct {
	UserID     uint64          `json:"user_id"`
	CategoryID uint64          `json:"category_id"`
	Year       uint            `json:"year"`
	Month      uint            `json:"month"`
	Budgeted   decimal.Decimal `json:"budgeted"`
	CreatedAt  time.Time       `json:"created_at"`
}

type UpdateBudgetRequest struct {
	CategoryID *uint64          `json:"category_id,omitempty"`
	Year       *uint            `json:"year,omitempty"`
	Month      *uint            `json:"month,omitempty"`
	Budgeted   *decimal.Decimal `json:"budgeted,omitempty"`
}

type CreateBudgetRequest struct {
	CategoryID uint64          `json:"category_id" binding:"required"`
	Year       uint            `json:"year" binding:"required"`
	Month      uint            `json:"month" binding:"required"`
	Budgeted   decimal.Decimal `json:"budgeted" binding:"required"`
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
