package repository_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"example.com/loan/module/loan/entity"
	"example.com/loan/module/loan/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRepaymentRepository_Create(t *testing.T) {
	testCases := []struct {
		name           string
		inpRepayment   entity.Repayment
		dbExecError    error
		lastInsertedID int64
		expectedErr    error
	}{
		{
			name:         "error: db connection error",
			inpRepayment: entity.Repayment{ID: "1", LoanBillingID: "1", Amount: decimal.NewFromInt(10000), Status: entity.RepaymentStatusCreated, Reference: "ref-1"},
			dbExecError:  errors.New("db connection error"),
			expectedErr:  errors.New("db connection error"),
		},
		{
			name:           "success: created",
			inpRepayment:   entity.Repayment{ID: "1", LoanBillingID: "1", Amount: decimal.NewFromInt(10000), Status: entity.RepaymentStatusCreated, Reference: "ref-1"},
			lastInsertedID: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)

			repo := repository.NewRepaymentRepository(db)
			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO repayments (loan_billing_id,amount,status,reference,created_at,updated_at) VALUES (?,?,?,?,?,?)")).
				WithArgs(
					tc.inpRepayment.LoanBillingID,
					tc.inpRepayment.Amount,
					tc.inpRepayment.Status,
					tc.inpRepayment.Reference,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).
				WillReturnResult(sqlmock.NewResult(tc.lastInsertedID, 1)).
				WillReturnError(tc.dbExecError)

			err = repo.Create(context.Background(), &tc.inpRepayment)
			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			assert.Equal(t, fmt.Sprint(tc.lastInsertedID), tc.inpRepayment.ID)
		})
	}
}

func TestRepaymentRepository_Update(t *testing.T) {
	testCases := []struct {
		name         string
		inpRepayment entity.Repayment
		dbExecError  error
		expectedErr  error
	}{
		{
			name:         "error: db connection error",
			inpRepayment: entity.Repayment{ID: "1", LoanBillingID: "1", Amount: decimal.NewFromInt(10000), Status: entity.RepaymentStatusPaid, Reference: "ref-1"},
			dbExecError:  errors.New("db connection error"),
			expectedErr:  errors.New("db connection error"),
		},
		{
			name:         "success: updated",
			inpRepayment: entity.Repayment{ID: "1", LoanBillingID: "1", Amount: decimal.NewFromInt(10000), Status: entity.RepaymentStatusSuccess, Reference: "ref-1"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)

			repo := repository.NewRepaymentRepository(db)
			mock.ExpectExec(regexp.QuoteMeta("UPDATE repayments SET amount = ?, status = ?, reference = ?, updated_at = ? WHERE id = ?")).
				WithArgs(
					tc.inpRepayment.Amount,
					tc.inpRepayment.Status,
					tc.inpRepayment.Reference,
					sqlmock.AnyArg(),
					tc.inpRepayment.ID,
				).
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tc.dbExecError)

			err = repo.Update(context.Background(), tc.inpRepayment)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
