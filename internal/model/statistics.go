package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type PeriodType string

const (
	PeriodTypeDay   PeriodType = "day"
	PeriodTypeWeek  PeriodType = "week"
	PeriodTypeMonth PeriodType = "month"
)

type CurrentBalanceStatistics struct {
	OnAccounts       decimal.Decimal `json:"on_accounts"`
	ReservedBills    decimal.Decimal `json:"reserved_bills"`
	ReservedBudgets  decimal.Decimal `json:"reserved_budgets"`
	TotalReserved    decimal.Decimal `json:"total_reserved"`
	FreeToDistribute decimal.Decimal `json:"free_to_distribute"`
}

type CategoryStatisticsRequest struct {
	Period PeriodType `json:"period" form:"period"`
	From   *time.Time `json:"from,omitempty" form:"from"`
	To     *time.Time `json:"to,omitempty" form:"to"`
}

type CategoryStatisticsItem struct {
	CategoryID   uint64          `json:"category_id" db:"category_id"`
	CategoryName string          `json:"category_name" db:"category_name"`
	Period       string          `json:"period" db:"period"`
	Income       decimal.Decimal `json:"income" db:"income"`
	Expense      decimal.Decimal `json:"expense" db:"expense"`
}

type CategoryStatisticsResponse struct {
	Period PeriodType               `json:"period"`
	Items  []CategoryStatisticsItem `json:"items"`
}

type PeriodStatisticsRequest struct {
	Period PeriodType `json:"period" form:"period"`
	From   *time.Time `json:"from,omitempty" form:"from"`
	To     *time.Time `json:"to,omitempty" form:"to"`
}

type PeriodStatisticsItem struct {
	Period  string          `json:"period" db:"period"`
	Income  decimal.Decimal `json:"income" db:"income"`
	Expense decimal.Decimal `json:"expense" db:"expense"`
	Balance decimal.Decimal `json:"balance" db:"balance"`
}

type PeriodStatisticsResponse struct {
	Period PeriodType             `json:"period"`
	Items  []PeriodStatisticsItem `json:"items"`
}
