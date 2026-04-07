package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"example.com/loan/module/payment/entity"
	sq "github.com/Masterminds/squirrel"
)

type mutationRepository struct {
	db *sql.DB
}

type MutationRepository interface {
	Get(ctx context.Context, mutation entity.Mutation) (entity.Mutation, error)
	Create(ctx context.Context, mutation *entity.Mutation) error
}

func (r *mutationRepository) Get(ctx context.Context, mutation entity.Mutation) (entity.Mutation, error) {
	query := sq.Select(
		"id",
		"account_id",
		"type",
		"reference",
		"amount",
	).From("mutations").Where(
		sq.Eq{
			"account_id": mutation.AccountID,
			"type":       mutation.Type,
			"reference":  mutation.Reference,
		},
	)

	row := query.RunWith(r.db).QueryRowContext(ctx)
	mutation, err := scanMutation(row)
	if err != nil {
		return entity.Mutation{}, err
	}

	return mutation, nil
}

func (r *mutationRepository) Create(ctx context.Context, mutation *entity.Mutation) error {
	gmt7 := time.FixedZone("GMT+7", 7*60*60)
	now := time.Now().In(gmt7)

	query := sq.Insert("mutations").Columns(
		"account_id",
		"type",
		"reference",
		"amount",
		"created_at",
		"updated_at",
	).Values(
		mutation.AccountID,
		mutation.Type,
		mutation.Reference,
		mutation.Amount,
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
	mutation.ID = fmt.Sprint(id)

	return nil
}

func scanMutation(row sq.RowScanner) (entity.Mutation, error) {
	var mutation entity.Mutation

	if err := row.Scan(
		&mutation.ID,
		&mutation.AccountID,
		&mutation.Type,
		&mutation.Reference,
		&mutation.Amount,
	); err != nil {
		if err == sql.ErrNoRows {
			return entity.Mutation{}, nil
		}

		return entity.Mutation{}, err
	}

	return mutation, nil
}

func NewMutationRepository(db *sql.DB) MutationRepository {
	return &mutationRepository{
		db: db,
	}
}
