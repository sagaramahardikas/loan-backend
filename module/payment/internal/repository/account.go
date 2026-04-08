package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"example.com/loan/module/payment/entity"
	sq "github.com/Masterminds/squirrel"
)

type accountRepository struct {
	db *sql.DB
}

type AccountRepository interface {
	GetByUserID(ctx context.Context, userID string) (entity.Account, error)
	Create(ctx context.Context, account *entity.Account) error
	Update(ctx context.Context, account entity.Account) error
}

func (r *accountRepository) GetByUserID(ctx context.Context, userID string) (entity.Account, error) {
	query := sq.Select(
		"id",
		"user_id",
		"balance",
		"status",
	).From("accounts").Where(sq.Eq{"user_id": userID})

	row := query.RunWith(r.db).QueryRowContext(ctx)
	account, err := scanAccount(row)
	if err != nil {
		return entity.Account{}, err
	}

	return account, nil
}

func (r *accountRepository) Update(ctx context.Context, account entity.Account) error {
	gmt7 := time.FixedZone("GMT+7", 7*60*60)
	now := time.Now().In(gmt7)

	query := sq.Update("accounts").
		Set("balance", account.Balance).
		Set("status", account.Status).
		Set("updated_at", now).
		Where(sq.Eq{"id": account.ID})

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

func (r *accountRepository) Create(ctx context.Context, account *entity.Account) error {
	gmt7 := time.FixedZone("GMT+7", 7*60*60)
	now := time.Now().In(gmt7)

	query := sq.Insert("accounts").Columns(
		"user_id",
		"balance",
		"status",
		"created_at",
		"updated_at",
	).Values(
		account.UserID,
		account.Balance,
		account.Status,
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
	account.ID = fmt.Sprint(id)

	return nil
}

func scanAccount(row sq.RowScanner) (entity.Account, error) {
	var account entity.Account

	if err := row.Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.Status,
	); err != nil {
		if err == sql.ErrNoRows {
			return entity.Account{}, nil
		}

		return entity.Account{}, err
	}

	return account, nil
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}
