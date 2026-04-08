package usecase_test

import (
	"context"
	"errors"
	"testing"

	"example.com/loan/module/payment/entity"
	"example.com/loan/module/payment/internal/repository/mock"
	"example.com/loan/module/payment/internal/usecase"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockAccountUsecase struct {
	repository *mock.MockAccountRepository
}

func TestAccountUsecase_Create(t *testing.T) {
	account := entity.Account{
		UserID:  "1",
		Balance: decimal.NewFromInt(10000),
		Status:  entity.AccountStatusActive,
	}

	testCases := []struct {
		name        string
		input       entity.Account
		mockFn      func(mock *mockAccountUsecase)
		expectedErr error
	}{
		{
			name:  "error: db connection error",
			input: account,
			mockFn: func(mocks *mockAccountUsecase) {
				mocks.repository.EXPECT().Create(
					gomock.Any(), &account,
				).Return(errors.New("db connection error"))
			},
			expectedErr: errors.New("db connection error"),
		},
		{
			name:  "success: created",
			input: account,
			mockFn: func(mocks *mockAccountUsecase) {
				mocks.repository.EXPECT().Create(
					gomock.Any(), &account,
				).Return(nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockAccountUsecase{repository: mock.NewMockAccountRepository(ctrl)}
			usecase := usecase.NewAccountUsecase(mock.repository)
			if tc.mockFn != nil {
				tc.mockFn(mock)
			}

			err := usecase.Create(context.Background(), &tc.input)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
