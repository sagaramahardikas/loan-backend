package entity

import (
	"github.com/shopspring/decimal"
)

//go:generate enumer -type=AccountStatus -trimprefix=AccountStatus -transform=kebab
type AccountStatus int8

const (
	AccountStatusUnspecified AccountStatus = iota
	AccountStatusInactive
	AccountStatusActive
	AccountStatusFrozen
)

type Account struct {
	ID      string          `json:"id"`
	UserID  string          `json:"user_id"`
	Balance decimal.Decimal `json:"balance"`
	Status  AccountStatus   `json:"status"`
}
