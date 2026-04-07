package entity

import "github.com/shopspring/decimal"

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
