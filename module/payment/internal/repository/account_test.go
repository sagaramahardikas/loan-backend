package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"example.com/loan/module/payment/entity"
	"example.com/loan/module/payment/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAccountRepository_GetByUserID(t *testing.T) {
	testCases := []struct {
		name            string
		userID          string
		dbGetError      error
		expectedAccount entity.Account
		expectedErr     error
	}{
		{
			name:        "error: not found",
			userID:      "1",
			dbGetError:  sql.ErrNoRows,
			expectedErr: sql.ErrNoRows,
		},
		{
			name:        "error: db connection error",
			userID:      "1",
			dbGetError:  errors.New("db connection error"),
			expectedErr: errors.New("db connection error"),
		},
		{
			name:            "success: found",
			userID:          "1",
			expectedAccount: entity.Account{ID: "1", UserID: "1", Balance: decimal.NewFromInt(10000), Status: entity.AccountStatusActive},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)
			repo := repository.NewAccountRepository(db)
			rows := mock.NewRows([]string{"id", "user_id", "balance", "status"})
			if tc.expectedErr == nil {
				rows.AddRow(
					tc.expectedAccount.ID,
					tc.expectedAccount.UserID,
					tc.expectedAccount.Balance,
					tc.expectedAccount.Status,
				)
			}

			mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, balance, status FROM accounts WHERE user_id = ?")).
				WithArgs(tc.userID).
				WillReturnRows(rows).
				WillReturnError(tc.dbGetError)

			got, err := repo.GetByUserID(context.Background(), tc.userID)
			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			assert.Equal(t, tc.expectedAccount, got)
		})
	}
}

func TestAccountRepository_Update(t *testing.T) {
	testCases := []struct {
		name        string
		inpAccount  entity.Account
		dbExecError error
		expectedErr error
	}{
		{
			name:        "error: db connection error",
			inpAccount:  entity.Account{ID: "1", UserID: "1", Balance: decimal.NewFromInt(10000), Status: entity.AccountStatusActive},
			dbExecError: errors.New("db connection error"),
			expectedErr: errors.New("db connection error"),
		},
		{
			name:       "success: updated",
			inpAccount: entity.Account{ID: "1", UserID: "1", Balance: decimal.NewFromInt(10000), Status: entity.AccountStatusActive},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)

			repo := repository.NewAccountRepository(db)
			mock.ExpectExec(regexp.QuoteMeta("UPDATE accounts SET balance = ?, status = ?, updated_at = ? WHERE id = ?")).
				WithArgs(
					tc.inpAccount.Balance,
					tc.inpAccount.Status,
					sqlmock.AnyArg(),
					tc.inpAccount.ID,
				).
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tc.dbExecError)

			err = repo.Update(context.Background(), tc.inpAccount)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
