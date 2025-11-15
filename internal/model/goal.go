package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type FrequencyType uint8

const (
	FrequencyTypeMonthly FrequencyType = iota
	FrequencyTypeDaily
	FrequencyTypeWeekly
	FrequencyTypeQuarterly
)

type Goal struct {
	ID           uint64          `json:"id" db:"id"`
	Name         string          `json:"name" db:"name"`
	TargetAmount decimal.Decimal `json:"target_amount" db:"target_amount"`
	StartAmount  decimal.Decimal `json:"start_amount" db:"start_amount"`
	Frequency    FrequencyType   `json:"frequency" db:"frequency"`
	DeadlineDate time.Time       `json:"deadline_date" db:"deadline_date"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
}
