package repository

import (
	"context"
	"database/sql"

	"example.com/loan/module/loan/entity"
	sq "github.com/Masterminds/squirrel"
)

type loanBillingRepository struct {
	db *sql.DB
}

type LoanBillingRepository interface {
	SumOutstandingLoans(ctx context.Context, loanID string) (entity.GetOutstandingLoansResponse, error)
}

func (r *loanBillingRepository) SumOutstandingLoans(ctx context.Context, loanID string) (entity.GetOutstandingLoansResponse, error) {
	query := sq.Select(
		"COALESCE(SUM(amount), 0) AS total_outstanding_amount",
		"COUNT(id) AS total_billing_count",
	).From("loan_billings").Where(
		sq.And{
			sq.Eq{"loan_id": loanID},
			sq.Lt{"status": entity.LoanBillingStatusPaid},
		},
	)

	var response entity.GetOutstandingLoansResponse
	if err := query.RunWith(r.db).QueryRowContext(ctx).Scan(
		&response.TotalOutstandingAmount,
		&response.TotalBillingCount,
	); err != nil {
		return entity.GetOutstandingLoansResponse{}, err
	}

	return response, nil
}

func NewLoanBillingRepository(db *sql.DB) LoanBillingRepository {
	return &loanBillingRepository{
		db: db,
	}
}
