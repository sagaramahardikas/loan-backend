package repository_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"example.com/loan/module/loan/entity"
	"example.com/loan/module/loan/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestLoanBillingRepository_GetByID(t *testing.T) {
	fixedDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name        string
		id          string
		dbGetError  error
		expectedObj entity.LoanBilling
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
			expectedObj: entity.LoanBilling{
				ID:      "123",
				LoanID:  "123",
				Amount:  decimal.NewFromInt(10000),
				Status:  entity.LoanBillingStatusPaid,
				DueDate: fixedDate,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)
			repo := repository.NewLoanBillingRepository(db)
			rows := mock.NewRows([]string{"id", "loan_id", "amount", "status", "due_date"})
			if tc.expectedErr == nil {
				rows.AddRow(
					tc.expectedObj.ID,
					tc.expectedObj.LoanID,
					tc.expectedObj.Amount,
					tc.expectedObj.Status,
					tc.expectedObj.DueDate,
				)
			}

			mock.ExpectQuery(regexp.QuoteMeta("SELECT id, loan_id, amount, status, due_date FROM loan_billings WHERE id = ?")).
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

func TestLoanBillingRepository_SumOutstandingLoans(t *testing.T) {
	testCases := []struct {
		name        string
		id          string
		dbListError error
		expectedObj entity.GetOutstandingLoansResponse
		expectedErr error
	}{
		{
			name:        "error: db connection error",
			id:          "123",
			dbListError: errors.New("db connection error"),
			expectedErr: errors.New("db connection error"),
		},
		{
			name: "success: found",
			id:   "123",
			expectedObj: entity.GetOutstandingLoansResponse{
				TotalOutstandingAmount: decimal.NewFromInt(100000),
				TotalBillingCount:      10,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)
			repo := repository.NewLoanBillingRepository(db)
			rows := mock.NewRows([]string{"total_outstanding_amount", "total_billing_count"})
			if tc.expectedErr == nil {
				rows.AddRow(
					tc.expectedObj.TotalOutstandingAmount,
					tc.expectedObj.TotalBillingCount,
				)
			}

			mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(amount), 0) AS total_outstanding_amount, COUNT(id) AS total_billing_count FROM loan_billings WHERE (loan_id = ? AND status <> ?)")).
				WithArgs(tc.id, entity.LoanBillingStatusPaid).
				WillReturnRows(rows).
				WillReturnError(tc.dbListError)

			got, err := repo.SumOutstandingLoans(context.Background(), tc.id)
			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			assert.Equal(t, tc.expectedObj, got)
		})
	}
}

func TestLoanBillingRepository_Update(t *testing.T) {
	testCases := []struct {
		name        string
		inpObj      entity.LoanBilling
		dbExecError error
		expectedErr error
	}{
		{
			name:        "error: db connection error",
			inpObj:      entity.LoanBilling{ID: "1", LoanID: "1", Amount: decimal.NewFromInt(10000), Status: entity.LoanBillingStatusPaid, DueDate: time.Now()},
			dbExecError: errors.New("db connection error"),
			expectedErr: errors.New("db connection error"),
		},
		{
			name:   "success: updated",
			inpObj: entity.LoanBilling{ID: "1", LoanID: "1", Amount: decimal.NewFromInt(10000), Status: entity.LoanBillingStatusPaid, DueDate: time.Now()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)

			repo := repository.NewLoanBillingRepository(db)
			mock.ExpectExec(regexp.QuoteMeta("UPDATE loan_billings SET amount = ?, status = ?, due_date = ?, updated_at = ? WHERE id = ?")).
				WithArgs(
					tc.inpObj.Amount,
					tc.inpObj.Status,
					tc.inpObj.DueDate,
					sqlmock.AnyArg(),
					tc.inpObj.ID,
				).
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tc.dbExecError)

			err = repo.Update(context.Background(), tc.inpObj)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestLoanBillingRepository_BulkCreate(t *testing.T) {
	fixedDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name        string
		inpBillings []entity.LoanBilling
		dbExecError error
		expectedErr error
	}{
		{
			name: "error: db connection error",
			inpBillings: []entity.LoanBilling{
				{
					LoanID:  "1",
					Amount:  decimal.NewFromInt(10000),
					Status:  entity.LoanBillingStatusCreated,
					DueDate: fixedDate.AddDate(0, 0, 7),
				},
				{
					LoanID:  "1",
					Amount:  decimal.NewFromInt(10000),
					Status:  entity.LoanBillingStatusCreated,
					DueDate: fixedDate.AddDate(0, 0, 14),
				},
			},
			dbExecError: errors.New("db connection error"),
			expectedErr: errors.New("db connection error"),
		},
		{
			name: "success: bulk create",
			inpBillings: []entity.LoanBilling{
				{
					LoanID:  "1",
					Amount:  decimal.NewFromInt(10000),
					Status:  entity.LoanBillingStatusCreated,
					DueDate: fixedDate.AddDate(0, 0, 7),
				},
				{
					LoanID:  "1",
					Amount:  decimal.NewFromInt(10000),
					Status:  entity.LoanBillingStatusCreated,
					DueDate: fixedDate.AddDate(0, 0, 14),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)

			repo := repository.NewLoanBillingRepository(db)
			billingsValues := ""
			argsValues := []driver.Value{}
			for i := range tc.inpBillings {
				billingsValues += "(?,?,?,?,?,?)"
				if i < len(tc.inpBillings)-1 {
					billingsValues += ","
				}

				argsValues = append(argsValues,
					tc.inpBillings[i].LoanID,
					tc.inpBillings[i].Amount,
					tc.inpBillings[i].Status,
					tc.inpBillings[i].DueDate,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				)
			}

			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO loan_billings (loan_id,amount,status,due_date,created_at,updated_at) VALUES " + billingsValues)).
				WithArgs(argsValues...).
				WillReturnResult(sqlmock.NewResult(0, 1)).
				WillReturnError(tc.dbExecError)

			err = repo.BulkCreate(context.Background(), tc.inpBillings)
			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}
		})
	}
}

func TestLoanBillingRepository_OverdueBillings(t *testing.T) {
	fixedDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name        string
		dbListError error
		expectedObj []entity.LoanBilling
		expectedErr error
	}{
		{
			name:        "error: db connection error",
			dbListError: errors.New("db connection error"),
			expectedErr: errors.New("db connection error"),
		},
		{
			name: "success: found all overdue billings",
			expectedObj: []entity.LoanBilling{
				{
					ID:      "1",
					LoanID:  "1",
					UserID:  "1",
					Amount:  decimal.NewFromInt(10000),
					Status:  entity.LoanBillingStatusCreated,
					DueDate: fixedDate.AddDate(0, 0, -7),
				},
				{
					ID:      "2",
					LoanID:  "1",
					UserID:  "1",
					Amount:  decimal.NewFromInt(10000),
					Status:  entity.LoanBillingStatusCreated,
					DueDate: fixedDate.AddDate(0, 0, -14),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.Nil(t, err)
			repo := repository.NewLoanBillingRepository(db)
			rows := mock.NewRows([]string{"id", "loan_id", "user_id", "amount", "status", "due_date"})
			if tc.expectedErr == nil {
				for _, billing := range tc.expectedObj {
					rows.AddRow(
						billing.ID,
						billing.LoanID,
						billing.UserID,
						billing.Amount,
						billing.Status,
						billing.DueDate,
					)
				}
			}

			mock.ExpectQuery(regexp.QuoteMeta("SELECT lb.id, lb.loan_id, l.user_id, lb.amount, lb.status, lb.due_date FROM loan_billings lb JOIN loans l ON lb.loan_id = l.id WHERE (lb.due_date < ? AND lb.status <> ?)")).
				WithArgs(sqlmock.AnyArg(), entity.LoanBillingStatusPaid).
				WillReturnRows(rows).
				WillReturnError(tc.dbListError)

			response, err := repo.OverdueBillings(context.Background())
			if err != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			assert.Equal(t, tc.expectedObj, response)
		})
	}
}
