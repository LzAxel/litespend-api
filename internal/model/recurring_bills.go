package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type RecurringBill struct {
	ID         uint64          `json:"id" db:"id"`
	UserID     uint64          `json:"user_id" db:"user_id"`
	Name       string          `json:"name" db:"name"`
	CategoryID *uint64         `json:"category_id,omitempty" db:"category_id"`
	Amount     decimal.Decimal `json:"amount" db:"amount"`
	DayDue     *int            `json:"day_due,omitempty" db:"day_due"`
	IsActive   bool            `json:"is_active" db:"is_active"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}

type BillInstance struct {
	ID              uint64          `json:"id" db:"id"`
	RecurringBillID uint64          `json:"recurring_bill_id" db:"recurring_bill_id"`
	Year            int             `json:"year" db:"year"`
	Month           int             `json:"month" db:"month"`
	AmountExpected  decimal.Decimal `json:"amount_expected" db:"amount_expected"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
}
