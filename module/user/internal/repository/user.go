package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"example.com/loan/module/user/entity"
	sq "github.com/Masterminds/squirrel"
)

type userRepository struct {
	db *sql.DB
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (entity.User, error)
	Create(ctx context.Context, user *entity.User) error
}

func (r *userRepository) GetByID(ctx context.Context, id string) (entity.User, error) {
	query := sq.Select(
		"id",
		"username",
		"status",
	).From("users").Where(sq.Eq{"id": id})

	row := query.RunWith(r.db).QueryRowContext(ctx)
	user, err := scanUser(row)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	gmt7 := time.FixedZone("GMT+7", 7*60*60)
	now := time.Now().In(gmt7)

	query := sq.Insert("users").Columns(
		"username",
		"status",
		"created_at",
		"updated_at",
	).Values(
		user.Username,
		user.Status,
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
	user.ID = fmt.Sprint(id)

	return nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func scanUser(row sq.RowScanner) (entity.User, error) {
	var user entity.User

	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Status,
	); err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, nil
		}

		return entity.User{}, err
	}

	return user, nil
}
