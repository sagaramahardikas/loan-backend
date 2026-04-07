package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"example.com/loan/module/loan/entity"
	sq "github.com/Masterminds/squirrel"
)

type repaymentRepository struct {
	db *sql.DB
}

type RepaymentRepository interface {
	Create(ctx context.Context, repayment *entity.Repayment) error
	Update(ctx context.Context, repayment entity.Repayment) error
}

func (r *repaymentRepository) Create(ctx context.Context, repayment *entity.Repayment) error {
	gmt7 := time.FixedZone("GMT+7", 7*60*60)
	now := time.Now().In(gmt7)

	query := sq.Insert("repayments").Columns(
		"loan_billing_id",
		"amount",
		"status",
		"reference",
		"created_at",
		"updated_at",
	).Values(
		repayment.LoanBillingID,
		repayment.Amount,
		repayment.Status,
		repayment.Reference,
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
	repayment.ID = fmt.Sprint(id)

	return nil
}

func (r *repaymentRepository) Update(ctx context.Context, repayment entity.Repayment) error {
	gmt7 := time.FixedZone("GMT+7", 7*60*60)
	now := time.Now().In(gmt7)

	query := sq.Update("repayments").
		Set("amount", repayment.Amount).
		Set("status", repayment.Status).
		Set("reference", repayment.Reference).
		Set("updated_at", now).
		Where(sq.Eq{"id": repayment.ID})

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

func NewRepaymentRepository(db *sql.DB) RepaymentRepository {
	return &repaymentRepository{
		db: db,
	}
}
