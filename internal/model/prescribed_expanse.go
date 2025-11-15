package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type PrescribedExpanse struct {
	ID          uint64          `json:"id" db:"id"`
	UserID      uint64          `json:"user_id" db:"user_id"`
	CategoryID  uint64          `json:"category_id" db:"category_id"`
	Description string          `json:"description" db:"description"`
	Frequency   FrequencyType   `json:"frequency" db:"frequency"`
	Amount      decimal.Decimal `json:"amount" db:"amount"`
	DateTime    time.Time       `json:"date_time" db:"date_time"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
}

type CreatePrescribedExpanseRequest struct {
	CategoryID  uint64          `json:"category_id" binding:"required"`
	Description string          `json:"description" binding:"required"`
	Frequency   FrequencyType   `json:"frequency" binding:"required"`
	Amount      decimal.Decimal `json:"amount" binding:"required"`
	DateTime    time.Time       `json:"date_time" binding:"required"`
}

type UpdatePrescribedExpanseRequest struct {
	CategoryID  *uint64          `json:"category_id"`
	Description *string          `json:"description"`
	Frequency   *FrequencyType   `json:"frequency"`
	Amount      *decimal.Decimal `json:"amount"`
	DateTime    *time.Time       `json:"date_time"`
}

type CreatePrescribedExpanseRecord struct {
	UserID      uint64
	CategoryID  uint64
	Description string
	Frequency   FrequencyType
	Amount      decimal.Decimal
	DateTime    time.Time
	CreatedAt   time.Time
}
