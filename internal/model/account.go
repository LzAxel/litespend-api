package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type AccountType string

const (
	AccountTypeCash   AccountType = "cash"
	AccountTypeBank   AccountType = "bank"
	AccountTypeCredit AccountType = "credit"
)

type Account struct {
	ID         uint64          `json:"id"`
	UserID     uint64          `json:"user_id"`
	Name       string          `json:"name"`
	Type       AccountType     `json:"type"`
	IsArchived bool            `json:"is_archived"`
	OrderNum   int             `json:"order_num"`
	Balance    decimal.Decimal `json:"balance"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type AccountDB struct {
	ID         uint64      `db:"id"`
	UserID     uint64      `db:"user_id"`
	Name       string      `db:"name"`
	Type       AccountType `db:"type"`
	IsArchived bool        `db:"is_archived"`
	OrderNum   int         `db:"order_num"`
	CreatedAt  time.Time   `db:"created_at"`
	UpdatedAt  time.Time   `db:"updated_at"`
}

type CreateAccountRequest struct {
	Name       string      `json:"name" db:"name"`
	Type       AccountType `json:"type" db:"type"`
	IsArchived bool        `json:"is_archived" db:"is_archived"`
	OrderNum   int         `json:"order_num" db:"order_num"`
}

type CreateAccountRecord struct {
	UserID     uint64
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
