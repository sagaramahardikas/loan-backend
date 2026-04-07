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

type mockMutationUsecase struct {
	accountRepo  *mock.MockAccountRepository
	mutationRepo *mock.MockMutationRepository
}

func TestMutationUsecase_Debit(t *testing.T) {
	//mutation := entity.Mutation{}

	testCases := []struct {
		name             string
		inpMutation      entity.DebitRequest
		mockFn           func(mock *mockMutationUsecase)
		expectedMutation entity.Mutation
		expectedErr      error
	}{
		{
			name: "error: account not found",
			inpMutation: entity.DebitRequest{
				UserID: "1",
			},
			mockFn: func(mocks *mockMutationUsecase) {
				mocks.accountRepo.EXPECT().GetByUserID(
					gomock.Any(), "1",
				).Return(entity.Account{}, errors.New("account not found"))
			},
			expectedErr: errors.New("account not found"),
		},
		{
			name: "error: account is not active",
			inpMutation: entity.DebitRequest{
				UserID: "1",
			},
			mockFn: func(mocks *mockMutationUsecase) {
				mocks.accountRepo.EXPECT().GetByUserID(
					gomock.Any(), "1",
				).Return(entity.Account{ID: "1", UserID: "1", Balance: decimal.NewFromInt(10000), Status: entity.AccountStatusInactive}, nil)
			},
			expectedErr: errors.New("account for user 1 is inactive"),
		},
		{
			name: "error: get mutation error",
			inpMutation: entity.DebitRequest{
				UserID:    "1",
				Type:      entity.MutationTypeRepayment,
				Reference: "ref-1",
			},
			mockFn: func(mocks *mockMutationUsecase) {
				mocks.accountRepo.EXPECT().GetByUserID(
					gomock.Any(), "1",
				).Return(entity.Account{ID: "1", UserID: "1", Balance: decimal.NewFromInt(10000), Status: entity.AccountStatusActive}, nil)

				mocks.mutationRepo.EXPECT().Get(
					gomock.Any(), entity.Mutation{AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1"},
				).Return(entity.Mutation{}, errors.New("get mutation error"))
			},
			expectedErr: errors.New("get mutation error"),
		},
		{
			name: "success: mutation already exist",
			inpMutation: entity.DebitRequest{
				UserID:    "1",
				Type:      entity.MutationTypeRepayment,
				Reference: "ref-1",
			},
			mockFn: func(mocks *mockMutationUsecase) {
				mocks.accountRepo.EXPECT().GetByUserID(
					gomock.Any(), "1",
				).Return(entity.Account{ID: "1", UserID: "1", Balance: decimal.NewFromInt(10000), Status: entity.AccountStatusActive}, nil)

				mocks.mutationRepo.EXPECT().Get(
					gomock.Any(), entity.Mutation{AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1"},
				).Return(entity.Mutation{ID: "1", AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1", Amount: decimal.NewFromInt(10000)}, nil)
			},
			expectedMutation: entity.Mutation{ID: "1", AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1", Amount: decimal.NewFromInt(10000)},
		},
		{
			name: "error: insufficient balance",
			inpMutation: entity.DebitRequest{
				UserID:    "1",
				Amount:    decimal.NewFromInt(20000),
				Type:      entity.MutationTypeRepayment,
				Reference: "ref-1",
			},
			mockFn: func(mocks *mockMutationUsecase) {
				mocks.accountRepo.EXPECT().GetByUserID(
					gomock.Any(), "1",
				).Return(entity.Account{ID: "1", UserID: "1", Balance: decimal.NewFromInt(10000), Status: entity.AccountStatusActive}, nil)

				mocks.mutationRepo.EXPECT().Get(
					gomock.Any(), entity.Mutation{AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1", Amount: decimal.NewFromInt(20000)},
				).Return(entity.Mutation{}, nil)
			},
			expectedErr: errors.New("insufficient balance for account 1"),
		},
		{
			name: "error: update account error",
			inpMutation: entity.DebitRequest{
				UserID:    "1",
				Amount:    decimal.NewFromInt(20000),
				Type:      entity.MutationTypeRepayment,
				Reference: "ref-1",
			},
			mockFn: func(mocks *mockMutationUsecase) {
				mocks.accountRepo.EXPECT().GetByUserID(
					gomock.Any(), "1",
				).Return(entity.Account{ID: "1", UserID: "1", Balance: decimal.NewFromInt(30000), Status: entity.AccountStatusActive}, nil)

				mocks.mutationRepo.EXPECT().Get(
					gomock.Any(), entity.Mutation{AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1", Amount: decimal.NewFromInt(20000)},
				).Return(entity.Mutation{}, nil)

				mocks.accountRepo.EXPECT().Update(
					gomock.Any(), entity.Account{ID: "1", UserID: "1", Balance: decimal.NewFromInt(10000), Status: entity.AccountStatusActive},
				).Return(errors.New("update account error"))
			},
			expectedErr: errors.New("update account error"),
		},
		{
			name: "error: create mutation error",
			inpMutation: entity.DebitRequest{
				UserID:    "1",
				Amount:    decimal.NewFromInt(20000),
				Type:      entity.MutationTypeRepayment,
				Reference: "ref-1",
			},
			mockFn: func(mocks *mockMutationUsecase) {
				mocks.accountRepo.EXPECT().GetByUserID(
					gomock.Any(), "1",
				).Return(entity.Account{ID: "1", UserID: "1", Balance: decimal.NewFromInt(30000), Status: entity.AccountStatusActive}, nil)

				mocks.mutationRepo.EXPECT().Get(
					gomock.Any(), entity.Mutation{AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1", Amount: decimal.NewFromInt(20000)},
				).Return(entity.Mutation{}, nil)

				mocks.accountRepo.EXPECT().Update(
					gomock.Any(), entity.Account{ID: "1", UserID: "1", Balance: decimal.NewFromInt(10000), Status: entity.AccountStatusActive},
				).Return(nil)

				mocks.mutationRepo.EXPECT().Create(
					gomock.Any(), &entity.Mutation{AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1", Amount: decimal.NewFromInt(20000)},
				).Return(errors.New("create mutation error"))
			},
			expectedErr: errors.New("create mutation error"),
		},
		{
			name: "success: debit",
			inpMutation: entity.DebitRequest{
				UserID:    "1",
				Amount:    decimal.NewFromInt(20000),
				Type:      entity.MutationTypeRepayment,
				Reference: "ref-1",
			},
			mockFn: func(mocks *mockMutationUsecase) {
				mocks.accountRepo.EXPECT().GetByUserID(
					gomock.Any(), "1",
				).Return(entity.Account{ID: "1", UserID: "1", Balance: decimal.NewFromInt(30000), Status: entity.AccountStatusActive}, nil)

				mocks.mutationRepo.EXPECT().Get(
					gomock.Any(), entity.Mutation{AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1", Amount: decimal.NewFromInt(20000)},
				).Return(entity.Mutation{}, nil)

				mocks.accountRepo.EXPECT().Update(
					gomock.Any(), entity.Account{ID: "1", UserID: "1", Balance: decimal.NewFromInt(10000), Status: entity.AccountStatusActive},
				).Return(nil)

				mocks.mutationRepo.EXPECT().Create(
					gomock.Any(), &entity.Mutation{AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1", Amount: decimal.NewFromInt(20000)},
				).Return(nil)
			},
			expectedMutation: entity.Mutation{AccountID: "1", Type: entity.MutationTypeRepayment, Reference: "ref-1", Amount: decimal.NewFromInt(20000)},
		},
	}

	ctrl := gomock.NewController(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockMutationUsecase{
				accountRepo:  mock.NewMockAccountRepository(ctrl),
				mutationRepo: mock.NewMockMutationRepository(ctrl),
			}

			usecase := usecase.NewMutationUsecase(mock.accountRepo, mock.mutationRepo)
			if tc.mockFn != nil {
				tc.mockFn(mock)
			}

			obj, err := usecase.Debit(context.Background(), tc.inpMutation)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedMutation, obj)
			}
		})
	}
}
