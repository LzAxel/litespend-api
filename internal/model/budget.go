package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type BudgetAllocation struct {
	ID         uint64          `json:"id" db:"id"`
	UserID     uint64          `json:"user_id" db:"user_id"`
	CategoryID uint64          `json:"category_id" db:"category_id"`
	Year       uint            `json:"year" db:"year"`
	Month      uint            `json:"month" db:"month"`
	Assigned   decimal.Decimal `json:"assigned" db:"assigned"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at" db:"updated_at"`
}

type CreateBudgetAllocationRecord struct {
	UserID     uint64          `json:"user_id"`
	CategoryID uint64          `json:"category_id"`
	Year       uint            `json:"year"`
	Month      uint            `json:"month"`
	Assigned   decimal.Decimal `json:"assigned"`
	CreatedAt  time.Time       `json:"created_at"`
}

type UpdateBudgetAllocationRequest struct {
	CategoryID *uint64          `json:"category_id,omitempty"`
	Year       *uint            `json:"year,omitempty"`
	Month      *uint            `json:"month,omitempty"`
	Assigned   *decimal.Decimal `json:"assigned,omitempty"`
}

type CreateBudgetAllocationRequest struct {
	CategoryID uint64          `json:"category_id" binding:"required"`
	Year       uint            `json:"year" binding:"required"`
	Month      uint            `json:"month" binding:"required"`
	Assigned   decimal.Decimal `json:"assigned" binding:"required"`
}
