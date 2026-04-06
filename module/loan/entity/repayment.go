package entity

import "github.com/shopspring/decimal"

//go:generate enumer -type=RepaymentStatus -trimprefix=RepaymentStatus -transform=kebab
type RepaymentStatus int8

const (
	RepaymentStatusUnspecified RepaymentStatus = iota
	RepaymentStatusCreated
	RepaymentStatusPaid
)

type Repayment struct {
	ID            string          `json:"id"`
	LoanBillingID string          `json:"loan_billing_id"`
	Amount        decimal.Decimal `json:"amount"`
	Status        RepaymentStatus `json:"status"`
	Reference     string          `json:"reference"`
}
