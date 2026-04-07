package repository

import (
	"context"
	"database/sql"
	"time"

	"example.com/loan/module/loan/entity"
	sq "github.com/Masterminds/squirrel"
)

type loanBillingRepository struct {
	db *sql.DB
}

type LoanBillingRepository interface {
	GetByID(ctx context.Context, billingID string) (entity.LoanBilling, error)
	SumOutstandingLoans(ctx context.Context, loanID string) (entity.GetOutstandingLoansResponse, error)
	Update(ctx context.Context, billing entity.LoanBilling) error
}

func (r *loanBillingRepository) GetByID(ctx context.Context, billingID string) (entity.LoanBilling, error) {
	query := sq.Select(
		"id",
		"loan_id",
		"amount",
		"status",
		"due_date",
	).From("loan_billings").Where(sq.Eq{"id": billingID})

	row := query.RunWith(r.db).QueryRowContext(ctx)
	billing, err := scanLoanBilling(row)
	if err != nil {
		return entity.LoanBilling{}, err
	}

	return billing, nil
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

func (r *loanBillingRepository) Update(ctx context.Context, billing entity.LoanBilling) error {
	gmt7 := time.FixedZone("GMT+7", 7*60*60)
	now := time.Now().In(gmt7)

	query := sq.Update("loan_billings").
		Set("amount", billing.Amount).
		Set("status", billing.Status).
		Set("due_date", billing.DueDate).
		Set("updated_at", now).
		Where(sq.Eq{"id": billing.ID})

	result, err := query.RunWith(r.db).ExecContext(ctx)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if affectedRows == 0 || err != nil {
		return sql.ErrNoRows
	}

	return err
}

func scanLoanBilling(row sq.RowScanner) (entity.LoanBilling, error) {
	var billing entity.LoanBilling

	if err := row.Scan(
		&billing.ID,
		&billing.LoanID,
		&billing.Amount,
		&billing.Status,
		&billing.DueDate,
	); err != nil {
		if err == sql.ErrNoRows {
			return entity.LoanBilling{}, nil
		}

		return entity.LoanBilling{}, err
	}

	return billing, nil
}

func NewLoanBillingRepository(db *sql.DB) LoanBillingRepository {
	return &loanBillingRepository{
		db: db,
	}
}
