package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Category struct {
	ID        uint64    `json:"id" db:"id"`
	UserID    uint64    `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	GroupName string    `json:"group_name" db:"group_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CategoryBudget struct {
	CategoryID  int64   `json:"category_id" db:"category_id"`
	Name        string  `json:"name" db:"category_name"`
	GroupName   string  `json:"group_name" db:"category_group_name"`
	Assigned    float64 `json:"assigned" db:"assigned"`
	Spent       float64 `json:"spent" db:"spent"`
	Available   float64 `json:"available" db:"available"`
	CarriedOver float64 `json:"carried_over" db:"carried_over"`
}

type CategoryBudgetResponse struct {
	ToBeBudgeted decimal.Decimal `json:"to_be_budgeted"`
	Categories []CategoryBudget `json:"categories"`
}

type CreateCategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	GroupName string `json:"group_name"`
}

type UpdateCategoryRequest struct {
	Name      *string `json:"name"`
	GroupName *string `json:"group_name"`
}

type UpdateCategoryRecord struct {
	Name      *string
	GroupName *string
	UpdatedAt time.Time
}

type CreateCategoryRecord struct {
	UserID    uint64
	Name      string
	GroupName string
	CreatedAt time.Time
	UpdatedAt time.Time
}
