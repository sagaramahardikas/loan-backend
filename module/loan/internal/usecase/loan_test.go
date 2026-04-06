package usecase_test

import (
	"context"
	"errors"
	"testing"

	"example.com/loan/module/loan/entity"
	"example.com/loan/module/loan/internal/repository/mock"
	"example.com/loan/module/loan/internal/usecase"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockLoanUsecase struct {
	usecase *mock.MockLoanBillingRepository
}

func TestLoanUsecase_GetOutstandingLoans(t *testing.T) {
	outstandingLoan := entity.GetOutstandingLoansResponse{
		TotalOutstandingAmount: decimal.NewFromInt(100000),
		TotalBillingCount:      10,
	}

	testCases := []struct {
		name        string
		id          string
		mockFn      func(mock *mockLoanUsecase)
		expectedObj entity.GetOutstandingLoansResponse
		expectedErr error
	}{
		{
			name: "error: db connection error",
			id:   "123",
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.usecase.EXPECT().SumOutstandingLoans(
					gomock.Any(), "123",
				).Return(entity.GetOutstandingLoansResponse{}, errors.New("db connection error"))
			},
			expectedErr: errors.New("db connection error"),
		},
		{
			name: "success: found",
			id:   "123",
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.usecase.EXPECT().SumOutstandingLoans(
					gomock.Any(), "123",
				).Return(outstandingLoan, nil)
			},
			expectedObj: outstandingLoan,
		},
	}

	ctrl := gomock.NewController(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockLoanUsecase{
				usecase: mock.NewMockLoanBillingRepository(ctrl),
			}

			usecase := usecase.NewLoanUsecase(mock.usecase)
			if tc.mockFn != nil {
				tc.mockFn(mock)
			}

			obj, err := usecase.GetOutstandingLoans(context.Background(), tc.id)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedObj, obj)
			}
		})
	}
}
