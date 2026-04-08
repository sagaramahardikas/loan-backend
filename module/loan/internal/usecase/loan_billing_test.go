package usecase_test

import (
	"context"
	"errors"
	"testing"

	"example.com/loan/module/loan/entity"
	"example.com/loan/module/loan/internal/repository/mock"
	"example.com/loan/module/loan/internal/usecase"
	"example.com/loan/module/user/client"
	mockUserClient "example.com/loan/module/user/client/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockLoanBillingUsecase struct {
	service     *mockUserClient.MockUserService
	billingRepo *mock.MockLoanBillingRepository
}

func TestLoanBillingUsecase_OverdueBillingChecker(t *testing.T) {
	billings := []entity.LoanBilling{
		{
			ID:     "1",
			UserID: "1",
			LoanID: "1",
			Amount: decimal.NewFromInt(10000),
			Status: entity.LoanBillingStatusCreated,
		},
		{
			ID:     "2",
			UserID: "1",
			LoanID: "1",
			Amount: decimal.NewFromInt(10000),
			Status: entity.LoanBillingStatusOverdue,
		},
	}

	updatedBilling := billings[0]
	updatedBilling.Status = entity.LoanBillingStatusOverdue

	testCases := []struct {
		name        string
		mockFn      func(mock *mockLoanBillingUsecase)
		expectedErr error
	}{
		{
			name: "error: failed to get overdue billings",
			mockFn: func(mocks *mockLoanBillingUsecase) {
				mocks.billingRepo.EXPECT().
					OverdueBillings(gomock.Any()).Return([]entity.LoanBilling{}, errors.New("get loan billing error"))
			},
			expectedErr: errors.New("get loan billing error"),
		},
		{
			name: "error: failed to update billing status",
			mockFn: func(mocks *mockLoanBillingUsecase) {
				mocks.billingRepo.EXPECT().
					OverdueBillings(gomock.Any()).Return(billings, nil)

				mocks.billingRepo.EXPECT().
					Update(gomock.Any(), updatedBilling).Return(errors.New("update billing error"))
			},
			expectedErr: errors.New("update billing error"),
		},
		{
			name: "error: failed to update user status",
			mockFn: func(mocks *mockLoanBillingUsecase) {
				mocks.billingRepo.EXPECT().
					OverdueBillings(gomock.Any()).Return(billings, nil)

				mocks.billingRepo.EXPECT().
					Update(gomock.Any(), updatedBilling).Return(nil)

				mocks.service.EXPECT().
					Update(gomock.Any(), client.UpdateUserRequest{
						UserID: "1",
						Status: 3,
					}).Return(client.UpdateUserResponse{}, errors.New("update user error"))
			},
			expectedErr: errors.New("update user error"),
		},
		{
			name: "success: check overdue billings and update user status",
			mockFn: func(mocks *mockLoanBillingUsecase) {
				mocks.billingRepo.EXPECT().
					OverdueBillings(gomock.Any()).Return(billings, nil)

				mocks.billingRepo.EXPECT().
					Update(gomock.Any(), updatedBilling).Return(nil)

				mocks.service.EXPECT().
					Update(gomock.Any(), client.UpdateUserRequest{
						UserID: "1",
						Status: 3,
					}).Return(client.UpdateUserResponse{}, nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockLoanBillingUsecase{
				service:     mockUserClient.NewMockUserService(ctrl),
				billingRepo: mock.NewMockLoanBillingRepository(ctrl),
			}

			usecase := usecase.NewLoanBillingUsecase(mock.billingRepo, mock.service)
			if tc.mockFn != nil {
				tc.mockFn(mock)
			}

			err := usecase.OverdueBillingChecker(context.Background())
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
