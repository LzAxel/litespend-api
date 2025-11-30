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

type CreateAccountRequest struct {
	Name       string      `json:"name" db:"name"`
	Type       AccountType `json:"type" db:"type"`
	IsArchived bool        `json:"is_archived" db:"is_archived"`
	OrderNum   int         `json:"order_num" db:"order_num"`
}

type CreateAccountRecord struct {
	UserID     int64
	Name       string
	Type       AccountType
	IsArchived bool
	OrderNum   int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UpdateAccountRequest struct {
	Name       *string `json:"name,omitempty"`
	IsArchived *bool   `json:"is_archived,omitempty"`
	OrderNum   *int    `json:"order_num,omitempty"`
}

type UpdateAccountRecord struct {
	Name       *string
	IsArchived *bool
	OrderNum   *int
	UpdatedAt  time.Time
}
