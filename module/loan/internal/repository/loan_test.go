package repository_test

import (
	"context"
	"database/sql"
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

func TestLoanRepository_GetByID(t *testing.T) {
	testCases := []struct {
		name        string
		id          string
		dbGetError  error
		expectedObj entity.Loan
		expectedErr error
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
			expectedObj: entity.Loan{
				ID:                "123",
				UserID:            "123",
				Principal:         decimal.NewFromInt(10000),
				Term:              10,
				Interest:          0.1,
				TotalAmount:       decimal.NewFromInt(11000),
				WeeklyInstallment: decimal.NewFromInt(1100),
				Status:            entity.LoanStatusProposed,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)
			repo := repository.NewLoanRepository(db)
			rows := mock.NewRows([]string{"id", "user_id", "principal", "term", "interest", "total_amount", "weekly_installment", "status"})
			if tc.expectedErr == nil {
				rows.AddRow(
					tc.expectedObj.ID,
					tc.expectedObj.UserID,
					tc.expectedObj.Principal,
					tc.expectedObj.Term,
					tc.expectedObj.Interest,
					tc.expectedObj.TotalAmount,
					tc.expectedObj.WeeklyInstallment,
					tc.expectedObj.Status,
				)
			}

			mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, principal, term, interest, total_amount, weekly_installment, status FROM loans WHERE id = ?")).
				WithArgs(tc.id).
				WillReturnRows(rows).
				WillReturnError(tc.dbGetError)

			got, err := repo.GetByID(context.Background(), tc.id)
			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			assert.Equal(t, tc.expectedObj, got)
		})
	}
}

func TestLoanRepository_Create(t *testing.T) {
	testCases := []struct {
		name           string
		inpLoan        entity.Loan
		dbExecError    error
		lastInsertedID int64
		expectedErr    error
	}{
		{
			name: "error: db connection error",
			inpLoan: entity.Loan{
				ID:                "1",
				UserID:            "1",
				Principal:         decimal.NewFromInt(10000),
				Term:              10,
				Interest:          0.1,
				TotalAmount:       decimal.NewFromInt(11000),
				WeeklyInstallment: decimal.NewFromInt(1100),
				Status:            entity.LoanStatusProposed,
			},
			dbExecError: errors.New("db connection error"),
			expectedErr: errors.New("db connection error"),
		},
		{
			name: "success: created",
			inpLoan: entity.Loan{
				ID:                "1",
				UserID:            "1",
				Principal:         decimal.NewFromInt(10000),
				Term:              10,
				Interest:          0.1,
				TotalAmount:       decimal.NewFromInt(11000),
				WeeklyInstallment: decimal.NewFromInt(1100),
				Status:            entity.LoanStatusProposed,
			},
			lastInsertedID: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)

			repo := repository.NewLoanRepository(db)
			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO loans (user_id,principal,term,interest,total_amount,weekly_installment,status,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?)")).
				WithArgs(
					tc.inpLoan.UserID,
					tc.inpLoan.Principal,
					tc.inpLoan.Term,
					tc.inpLoan.Interest,
					tc.inpLoan.TotalAmount,
					tc.inpLoan.WeeklyInstallment,
					tc.inpLoan.Status,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).
				WillReturnResult(sqlmock.NewResult(tc.lastInsertedID, 1)).
				WillReturnError(tc.dbExecError)

			err = repo.Create(context.Background(), &tc.inpLoan)
			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			assert.Equal(t, fmt.Sprint(tc.lastInsertedID), tc.inpLoan.ID)
		})
	}
}

func TestLoanRepository_Update(t *testing.T) {
	testCases := []struct {
		name        string
		inpLoan     entity.Loan
		dbExecError error
		expectedErr error
	}{
		{
			name:        "error: db connection error",
			inpLoan:     entity.Loan{ID: "1", UserID: "1", Principal: decimal.NewFromInt(10000), Term: 10, Interest: 0.1, TotalAmount: decimal.NewFromInt(11000), WeeklyInstallment: decimal.NewFromInt(1100), Status: entity.LoanStatusProposed},
			dbExecError: errors.New("db connection error"),
			expectedErr: errors.New("db connection error"),
		},
		{
			name:    "success: updated",
			inpLoan: entity.Loan{ID: "1", UserID: "1", Principal: decimal.NewFromInt(10000), Term: 10, Interest: 0.1, TotalAmount: decimal.NewFromInt(11000), WeeklyInstallment: decimal.NewFromInt(1100), Status: entity.LoanStatusProposed},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)

			repo := repository.NewLoanRepository(db)
			mock.ExpectExec(regexp.QuoteMeta("UPDATE loans SET principal = ?, term = ?, interest = ?, total_amount = ?, weekly_installment = ?, status = ?, updated_at = ? WHERE id = ?")).
				WithArgs(
					tc.inpLoan.Principal,
					tc.inpLoan.Term,
					tc.inpLoan.Interest,
					tc.inpLoan.TotalAmount,
					tc.inpLoan.WeeklyInstallment,
					tc.inpLoan.Status,
					sqlmock.AnyArg(),
					tc.inpLoan.ID,
				).
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tc.dbExecError)

			err = repo.Update(context.Background(), tc.inpLoan)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
