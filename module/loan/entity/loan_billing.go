package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type GetOutstandingLoansResponse struct {
	TotalOutstandingAmount decimal.Decimal `json:"total_outstanding_amount"`
	TotalBillingCount      int             `json:"total_billing_count"`
}

//go:generate enumer -type=LoanBillingStatus -trimprefix=LoanBillingStatus -transform=kebab
type LoanBillingStatus int8

const (
	LoanBillingStatusUnspecified LoanBillingStatus = iota
	LoanBillingStatusCreated
	LoanBillingStatusPaid
	LoanBillingStatusOverdue
)

type LoanBilling struct {
	ID      string            `json:"id"`
	LoanID  string            `json:"loan_id"`
	Amount  decimal.Decimal   `json:"amount"`
	Status  LoanBillingStatus `json:"status"`
	DueDate time.Time         `json:"due_date"`
}
