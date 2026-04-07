package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"example.com/loan/module/payment/entity"
	"example.com/loan/module/payment/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestMutationRepository_Get(t *testing.T) {
	testCases := []struct {
		name             string
		inpMutation      entity.Mutation
		dbGetError       error
		expectedMutation entity.Mutation
		expectedErr      error
	}{
		{
			name:        "error: not found",
			inpMutation: entity.Mutation{AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1"},
			dbGetError:  sql.ErrNoRows,
			expectedErr: sql.ErrNoRows,
		},
		{
			name:        "error: db connection error",
			inpMutation: entity.Mutation{AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1"},
			dbGetError:  errors.New("db connection error"),
			expectedErr: errors.New("db connection error"),
		},
		{
			name:             "success: found",
			inpMutation:      entity.Mutation{AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1"},
			expectedMutation: entity.Mutation{ID: "1", AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1", Amount: decimal.NewFromInt(10000)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)
			repo := repository.NewMutationRepository(db)
			rows := mock.NewRows([]string{"id", "account_id", "type", "reference", "amount"})
			if tc.expectedErr == nil {
				rows.AddRow(
					tc.expectedMutation.ID,
					tc.expectedMutation.AccountID,
					tc.expectedMutation.Type,
					tc.expectedMutation.Reference,
					tc.expectedMutation.Amount,
				)
			}

			mock.ExpectQuery(regexp.QuoteMeta("SELECT id, account_id, type, reference, amount FROM mutations WHERE account_id = ? AND reference = ? AND type = ?")).
				WithArgs(tc.inpMutation.AccountID, tc.inpMutation.Reference, tc.inpMutation.Type).
				WillReturnRows(rows).
				WillReturnError(tc.dbGetError)

			got, err := repo.Get(context.Background(), tc.inpMutation)
			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			assert.Equal(t, tc.expectedMutation, got)
		})
	}
}

func TestMutationRepository_Create(t *testing.T) {
	testCases := []struct {
		name           string
		inpMutation    entity.Mutation
		dbExecError    error
		lastInsertedID int64
		expectedErr    error
	}{
		{
			name:        "error: db connection error",
			inpMutation: entity.Mutation{AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1", Amount: decimal.NewFromInt(10000)},
			dbExecError: errors.New("db connection error"),
			expectedErr: errors.New("db connection error"),
		},
		{
			name:           "success: created",
			inpMutation:    entity.Mutation{AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1", Amount: decimal.NewFromInt(10000)},
			lastInsertedID: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)

			repo := repository.NewMutationRepository(db)
			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO mutations (account_id,type,reference,amount,created_at,updated_at) VALUES (?,?,?,?,?,?)")).
				WithArgs(
					tc.inpMutation.AccountID,
					tc.inpMutation.Type,
					tc.inpMutation.Reference,
					tc.inpMutation.Amount,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).
				WillReturnResult(sqlmock.NewResult(tc.lastInsertedID, 1)).
				WillReturnError(tc.dbExecError)

			err = repo.Create(context.Background(), &tc.inpMutation)
			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			assert.Equal(t, fmt.Sprint(tc.lastInsertedID), tc.inpMutation.ID)
		})
	}
}
