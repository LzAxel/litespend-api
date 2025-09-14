package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type GoalFrequencyType uint8

const (
	GoalFrequencyTypeMonthly GoalFrequencyType = iota
	GoalFrequencyTypeDaily
	GoalFrequencyTypeWeekly
	GoalFrequencyTypeQuarterly
)

type Goal struct {
	ID           int               `json:"id" db:"id"`
	UserID       int               `json:"user_id" db:"user_id"`
	Name         string            `json:"name" db:"name"`
	TargetAmount decimal.Decimal   `json:"target_amount" db:"target_amount"`
	StartAmount  decimal.Decimal   `json:"start_amount" db:"start_amount"`
	Frequency    GoalFrequencyType `json:"frequency" db:"frequency"`
	DeadlineDate time.Time         `json:"deadline_date" db:"deadline_date"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
}

type CreateGoalRecord struct {
	UserID       int
	Name         string
	TargetAmount decimal.Decimal
	StartAmount  decimal.Decimal
	Frequency    GoalFrequencyType
	DeadlineDate time.Time
	CreatedAt    time.Time
}

type UpdateGoalRecord struct {
	UserID       int
	Name         string
	TargetAmount decimal.Decimal
	StartAmount  decimal.Decimal
	Frequency    GoalFrequencyType
	DeadlineDate time.Time
}
