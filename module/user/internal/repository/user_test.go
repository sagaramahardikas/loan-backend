package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"example.com/loan/module/user/entity"
	"example.com/loan/module/user/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_GetByID(t *testing.T) {
	testCases := []struct {
		name         string
		id           string
		dbGetError   error
		expectedUser entity.User
		expectedErr  error
	}{
		{
			name:        "error: not found",
			id:          "123",
			dbGetError:  sql.ErrNoRows,
			expectedErr: sql.ErrNoRows,
		},
		{
			name:        "error: db connection error",
			id:          "123",
			dbGetError:  errors.New("db connection error"),
			expectedErr: errors.New("db connection error"),
		},
		{
			name: "success: found",
			id:   "123",
			expectedUser: entity.User{
				ID:       "123",
				Username: "testuser",
				Status:   entity.UserStatusActive,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)
			repo := repository.NewUserRepository(db)
			rows := mock.NewRows([]string{"id", "username", "status"})
			if tc.expectedErr == nil {
				rows.AddRow(
					tc.expectedUser.ID,
					tc.expectedUser.Username,
					tc.expectedUser.Status,
				)
			}

			mock.ExpectQuery(regexp.QuoteMeta("SELECT id, username, status FROM users WHERE id = ?")).
				WithArgs(tc.id).
				WillReturnRows(rows).
				WillReturnError(tc.dbGetError)

			got, err := repo.GetByID(context.Background(), tc.id)
			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			assert.Equal(t, tc.expectedUser, got)
		})
	}
}
