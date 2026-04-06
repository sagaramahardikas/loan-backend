package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"example.com/loan/module/loan/entity"
	"example.com/loan/module/loan/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

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

			mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(amount), 0) AS total_outstanding_amount, COUNT(id) AS total_billing_count FROM loan_billings WHERE (loan_id = ? AND status < ?)")).
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
