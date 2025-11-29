package model

import "time"

type AccountType string

const (
	AccountTypeCash   AccountType = "cash"
	AccountTypeBank   AccountType = "bank"
	AccountTypeCredit AccountType = "credit"
)

type Account struct {
	ID         int64       `json:"id" db:"id"`
	UserID     int64       `json:"user_id" db:"user_id"`
	Name       string      `json:"name" db:"name"`
	Type       AccountType `json:"type" db:"type"`
	IsArchived bool        `json:"is_archived" db:"is_archived"`
	OrderNum   int         `json:"order_num" db:"order_num"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at" db:"updated_at"`
}
