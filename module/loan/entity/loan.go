package entity

import (
	"github.com/shopspring/decimal"
)

//go:generate enumer -type=LoanStatus -trimprefix=LoanStatus -transform=kebab
type LoanStatus int8

const (
	LoanStatusUnspecified LoanStatus = iota
	LoanStatusProposed
	LoanStatusApproved
	LoanStatusInvested
	LoanStatusDisbursed
)

type Loan struct {
	ID                string          `json:"id"`
	UserID            string          `json:"user_id"`
	Principal         decimal.Decimal `json:"principal"`
	Term              int             `json:"term"`     // in months
	Interest          float64         `json:"interest"` // annual interest rate
	TotalAmount       decimal.Decimal `json:"total_amount"`
	WeeklyInstallment decimal.Decimal `json:"weekly_installment"`
	Status            LoanStatus      `json:"status"`
}
