package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type GetOutstandingLoansResponse struct {
	TotalOutstandingAmount decimal.Decimal `json:"total_outstanding_amount"`
	TotalBillingCount      int             `json:"total_billing_count"`
}

type PayBillingRequest struct {
	UserID    string          `json:"user_id"`
	BillingID string          `json:"billing_id"`
	Amount    decimal.Decimal `json:"amount"`
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
	UserID  string            `json:"user_id,omitempty"` // joined from loans
	LoanID  string            `json:"loan_id"`
	Amount  decimal.Decimal   `json:"amount"`
	Status  LoanBillingStatus `json:"status"`
	DueDate time.Time         `json:"due_date"`
}
