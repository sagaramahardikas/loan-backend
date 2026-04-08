package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"example.com/loan/module/loan/entity"
	sq "github.com/Masterminds/squirrel"
)

type loanRepository struct {
	db *sql.DB
}

type LoanRepository interface {
	GetByID(ctx context.Context, loanID string) (entity.Loan, error)
	Create(ctx context.Context, loan *entity.Loan) error
	Update(ctx context.Context, loan entity.Loan) error
}

func (r *loanRepository) GetByID(ctx context.Context, loanID string) (entity.Loan, error) {
	query := sq.Select(
		"id",
		"user_id",
		"principal",
		"term",
		"interest",
		"total_amount",
		"weekly_installment",
		"status",
	).From("loans").Where(sq.Eq{"id": loanID})

	row := query.RunWith(r.db).QueryRowContext(ctx)
	loan, err := scanLoan(row)
	if err != nil {
		return entity.Loan{}, err
	}

	return loan, nil
}

func (r *loanRepository) Create(ctx context.Context, loan *entity.Loan) error {
	gmt7 := time.FixedZone("GMT+7", 7*60*60)
	now := time.Now().In(gmt7)

	query := sq.Insert("loans").Columns(
		"user_id",
		"principal",
		"term",
		"interest",
		"total_amount",
		"weekly_installment",
		"status",
		"created_at",
		"updated_at",
	).Values(
		loan.UserID,
		loan.Principal,
		loan.Term,
		loan.Interest,
		loan.TotalAmount,
		loan.WeeklyInstallment,
		loan.Status,
		now,
		now,
	)

	result, err := query.RunWith(r.db).ExecContext(ctx)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	loan.ID = fmt.Sprint(id)

	return nil
}

func (r *loanRepository) Update(ctx context.Context, loan entity.Loan) error {
	gmt7 := time.FixedZone("GMT+7", 7*60*60)
	now := time.Now().In(gmt7)

	query := sq.Update("loans").
		Set("principal", loan.Principal).
		Set("term", loan.Term).
		Set("interest", loan.Interest).
		Set("total_amount", loan.TotalAmount).
		Set("weekly_installment", loan.WeeklyInstallment).
		Set("status", loan.Status).
		Set("updated_at", now).
		Where(sq.Eq{"id": loan.ID})

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

func scanLoan(row sq.RowScanner) (entity.Loan, error) {
	var loan entity.Loan
	if err := row.Scan(
		&loan.ID,
		&loan.UserID,
		&loan.Principal,
		&loan.Term,
		&loan.Interest,
		&loan.TotalAmount,
		&loan.WeeklyInstallment,
		&loan.Status,
	); err != nil {
		if err == sql.ErrNoRows {
			return entity.Loan{}, nil
		}

		return entity.Loan{}, err
	}

	return loan, nil
}

func NewLoanRepository(db *sql.DB) LoanRepository {
	return &loanRepository{
		db: db,
	}
}
