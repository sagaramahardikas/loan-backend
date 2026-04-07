package entity

import "github.com/shopspring/decimal"

type DebitRequest struct {
	UserID    string          `json:"user_id"`
	Amount    decimal.Decimal `json:"amount"`
	Type      MutationType    `json:"type"`
	Reference string          `json:"reference"`
}

//go:generate enumer -type=MutationType -trimprefix=MutationType -transform=kebab
type MutationType int8

const (
	MutationTypeUnspecified MutationType = iota
	MutationTypeRepayment
)

type Mutation struct {
	ID        string          `json:"id"`
	AccountID string          `json:"account_id"`
	Type      MutationType    `json:"type"`
	Reference string          `json:"reference"`
	Amount    decimal.Decimal `json:"amount"`
}
