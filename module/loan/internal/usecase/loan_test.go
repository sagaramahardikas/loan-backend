package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"example.com/loan/internal/util"
	"example.com/loan/module/loan/entity"
	"example.com/loan/module/loan/internal/repository/mock"
	"example.com/loan/module/loan/internal/usecase"
	"example.com/loan/module/payment/client"
	mockPaymentClient "example.com/loan/module/payment/client/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockLoanUsecase struct {
	service       *mockPaymentClient.MockPaymentService
	billingRepo   *mock.MockLoanBillingRepository
	repaymentRepo *mock.MockRepaymentRepository
	loanRepo      *mock.MockLoanRepository
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
				mocks.billingRepo.EXPECT().SumOutstandingLoans(
					gomock.Any(), "123",
				).Return(entity.GetOutstandingLoansResponse{}, errors.New("db connection error"))
			},
			expectedErr: errors.New("db connection error"),
		},
		{
			name: "success: found",
			id:   "123",
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.billingRepo.EXPECT().SumOutstandingLoans(
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
				service:       mockPaymentClient.NewMockPaymentService(ctrl),
				billingRepo:   mock.NewMockLoanBillingRepository(ctrl),
				repaymentRepo: mock.NewMockRepaymentRepository(ctrl),
				loanRepo:      mock.NewMockLoanRepository(ctrl),
			}

			usecase := usecase.NewLoanUsecase(mock.service, mock.billingRepo, mock.repaymentRepo, mock.loanRepo, nil)
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

func TestLoanUsecase_PayBilling(t *testing.T) {
	fixedDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)

	loanBillingCreated := entity.LoanBilling{
		ID:      "123",
		LoanID:  "123",
		Amount:  decimal.NewFromInt(10000),
		Status:  entity.LoanBillingStatusCreated,
		DueDate: fixedDate,
	}

	loanBillingPaid := loanBillingCreated
	loanBillingPaid.Status = entity.LoanBillingStatusPaid

	repayment := entity.Repayment{
		LoanBillingID: "123",
		Amount:        decimal.NewFromInt(10000),
		Status:        entity.RepaymentStatusCreated,
	}

	request := entity.PayBillingRequest{
		UserID:    "1",
		BillingID: "123",
		Amount:    decimal.NewFromInt(10000),
	}

	createAndPayMutationReq := client.CreateAndPayMutationRequest{
		UserID:    "1",
		Amount:    decimal.NewFromInt(10000),
		Reference: fmt.Sprintf("%s-%s", usecase.RepaymentReferencePrefix, "1"),
	}

	repaymentUpdate := repayment
	repaymentUpdate.ID = "1"
	repaymentUpdate.Reference = createAndPayMutationReq.Reference

	testCases := []struct {
		name        string
		input       entity.PayBillingRequest
		mockFn      func(mock *mockLoanUsecase)
		expectedObj entity.LoanBilling
		expectedErr error
	}{
		{
			name:  "error: get loan billing error",
			input: request,
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.billingRepo.EXPECT().GetByID(
					gomock.Any(), "123",
				).Return(entity.LoanBilling{}, errors.New("get loan billing error"))
			},
			expectedErr: errors.New("get loan billing error"),
		},
		{
			name:  "success: loan billing already paid",
			input: request,
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.billingRepo.EXPECT().GetByID(
					gomock.Any(), "123",
				).Return(loanBillingPaid, nil)
			},
			expectedObj: loanBillingPaid,
		},
		{
			name:  "error: repayment create error",
			input: request,
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.billingRepo.EXPECT().GetByID(
					gomock.Any(), "123",
				).Return(loanBillingCreated, nil)

				mocks.repaymentRepo.EXPECT().Create(
					gomock.Any(), &repayment,
				).Return(errors.New("repayment create error"))
			},
			expectedErr: errors.New("repayment create error"),
		},
		{
			name:  "error: payment service error",
			input: request,
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.billingRepo.EXPECT().GetByID(
					gomock.Any(), "123",
				).Return(loanBillingCreated, nil)

				mocks.repaymentRepo.EXPECT().Create(
					gomock.Any(), &repayment,
				).Do(func(ctx context.Context, repayment *entity.Repayment) {
					repayment.ID = "1"
				}).Return(nil)

				mocks.service.EXPECT().CreateAndPayMutation(
					gomock.Any(), createAndPayMutationReq,
				).Return(client.CreateAndPayMutationResponse{}, errors.New("payment service error"))
			},
			expectedErr: errors.New("payment service error"),
		},
		{
			name:  "error: loan billing update error",
			input: request,
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.billingRepo.EXPECT().GetByID(
					gomock.Any(), "123",
				).Return(loanBillingCreated, nil)

				mocks.repaymentRepo.EXPECT().Create(
					gomock.Any(), &repayment,
				).Do(func(ctx context.Context, repayment *entity.Repayment) {
					repayment.ID = "1"
				}).Return(nil)

				mocks.service.EXPECT().CreateAndPayMutation(
					gomock.Any(), createAndPayMutationReq,
				).Return(client.CreateAndPayMutationResponse{}, nil)

				mocks.billingRepo.EXPECT().Update(
					gomock.Any(), loanBillingPaid,
				).Return(errors.New("loan billing update error"))
			},
			expectedErr: errors.New("loan billing update error"),
		},
		{
			name:  "error: repayment update error",
			input: request,
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.billingRepo.EXPECT().GetByID(
					gomock.Any(), "123",
				).Return(loanBillingCreated, nil)

				mocks.repaymentRepo.EXPECT().Create(
					gomock.Any(), &repayment,
				).Do(func(ctx context.Context, repayment *entity.Repayment) {
					repayment.ID = "1"
				}).Return(nil)

				mocks.service.EXPECT().CreateAndPayMutation(
					gomock.Any(), createAndPayMutationReq,
				).Return(client.CreateAndPayMutationResponse{}, nil)

				mocks.billingRepo.EXPECT().Update(
					gomock.Any(), loanBillingPaid,
				).Return(nil)

				mocks.repaymentRepo.EXPECT().Update(
					gomock.Any(), repaymentUpdate,
				).Return(errors.New("repayment update error"))
			},
			expectedErr: errors.New("repayment update error"),
		},
		{
			name:  "success: pay billing",
			input: request,
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.billingRepo.EXPECT().GetByID(
					gomock.Any(), "123",
				).Return(loanBillingCreated, nil)

				mocks.repaymentRepo.EXPECT().Create(
					gomock.Any(), &repayment,
				).Do(func(ctx context.Context, repayment *entity.Repayment) {
					repayment.ID = "1"
				}).Return(nil)

				mocks.service.EXPECT().CreateAndPayMutation(
					gomock.Any(), createAndPayMutationReq,
				).Return(client.CreateAndPayMutationResponse{}, nil)

				mocks.billingRepo.EXPECT().Update(
					gomock.Any(), loanBillingPaid,
				).Return(nil)

				mocks.repaymentRepo.EXPECT().Update(
					gomock.Any(), repaymentUpdate,
				).Return(nil)
			},
			expectedObj: loanBillingPaid,
		},
	}

	ctrl := gomock.NewController(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockLoanUsecase{
				service:       mockPaymentClient.NewMockPaymentService(ctrl),
				billingRepo:   mock.NewMockLoanBillingRepository(ctrl),
				repaymentRepo: mock.NewMockRepaymentRepository(ctrl),
				loanRepo:      mock.NewMockLoanRepository(ctrl),
			}

			usecase := usecase.NewLoanUsecase(mock.service, mock.billingRepo, mock.repaymentRepo, mock.loanRepo, nil)
			if tc.mockFn != nil {
				tc.mockFn(mock)
			}

			obj, err := usecase.PayBilling(context.Background(), tc.input)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedObj, obj)
			}
		})
	}
}

func TestLoanUsecase_Create(t *testing.T) {
	loan := entity.Loan{
		UserID:            "1",
		Principal:         decimal.NewFromInt(100000),
		Term:              10,
		Interest:          0.1,
		TotalAmount:       decimal.NewFromInt(110000),
		WeeklyInstallment: decimal.NewFromInt(10000),
		Status:            entity.LoanStatusProposed,
	}

	testCases := []struct {
		name        string
		input       entity.Loan
		mockFn      func(mock *mockLoanUsecase)
		expectedErr error
	}{
		{
			name:  "error: db connection error",
			input: loan,
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.loanRepo.EXPECT().Create(
					gomock.Any(), &loan,
				).Return(errors.New("db connection error"))
			},
			expectedErr: errors.New("db connection error"),
		},
		{
			name:  "success: created",
			input: loan,
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.loanRepo.EXPECT().Create(
					gomock.Any(), &loan,
				).Return(nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockLoanUsecase{loanRepo: mock.NewMockLoanRepository(ctrl)}
			usecase := usecase.NewLoanUsecase(nil, nil, nil, mock.loanRepo, nil)
			if tc.mockFn != nil {
				tc.mockFn(mock)
			}

			err := usecase.Create(context.Background(), &tc.input)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestLoanUsecase_ForceDisburse(t *testing.T) {
	fixedTime := time.Date(2026, 4, 8, 0, 0, 0, 0, time.FixedZone("GMT+7", 7*60*60))
	clock := util.FixedClock{FixedTime: fixedTime}

	loan := entity.Loan{
		UserID:            "1",
		Principal:         decimal.NewFromInt(100000),
		Term:              2,
		Interest:          0.1,
		TotalAmount:       decimal.NewFromInt(110000),
		WeeklyInstallment: decimal.NewFromInt(10000),
		Status:            entity.LoanStatusProposed,
	}

	disbursedLoan := loan
	disbursedLoan.Status = entity.LoanStatusDisbursed

	billings := []entity.LoanBilling{
		{
			LoanID:  loan.ID,
			Amount:  loan.WeeklyInstallment,
			Status:  entity.LoanBillingStatusCreated,
			DueDate: fixedTime.AddDate(0, 0, 7),
		},
		{
			LoanID:  loan.ID,
			Amount:  loan.WeeklyInstallment,
			Status:  entity.LoanBillingStatusCreated,
			DueDate: fixedTime.AddDate(0, 0, 14),
		},
	}

	testCases := []struct {
		name        string
		loanID      string
		mockFn      func(mock *mockLoanUsecase)
		expectedErr error
	}{
		{
			name:   "error: get loan error",
			loanID: loan.ID,
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.loanRepo.EXPECT().GetByID(
					gomock.Any(), loan.ID,
				).Return(entity.Loan{}, errors.New("get loan error"))
			},
			expectedErr: errors.New("get loan error"),
		},
		{
			name:   "error: update loan error",
			loanID: loan.ID,
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.loanRepo.EXPECT().GetByID(
					gomock.Any(), loan.ID,
				).Return(loan, nil)

				mocks.loanRepo.EXPECT().Update(
					gomock.Any(), disbursedLoan,
				).Return(errors.New("update loan error"))
			},
			expectedErr: errors.New("update loan error"),
		},
		{
			name:   "error: loan bulk create error",
			loanID: loan.ID,
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.loanRepo.EXPECT().GetByID(
					gomock.Any(), loan.ID,
				).Return(loan, nil)

				mocks.loanRepo.EXPECT().Update(
					gomock.Any(), disbursedLoan,
				).Return(nil)

				mocks.billingRepo.EXPECT().BulkCreate(
					gomock.Any(), billings,
				).Return(errors.New("bulk create loan billing error"))
			},
			expectedErr: errors.New("bulk create loan billing error"),
		},
		{
			name:   "success: bulk create",
			loanID: loan.ID,
			mockFn: func(mocks *mockLoanUsecase) {
				mocks.loanRepo.EXPECT().GetByID(
					gomock.Any(), loan.ID,
				).Return(loan, nil)

				mocks.loanRepo.EXPECT().Update(
					gomock.Any(), disbursedLoan,
				).Return(nil)

				mocks.billingRepo.EXPECT().BulkCreate(
					gomock.Any(), billings,
				).Return(nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockLoanUsecase{
				billingRepo: mock.NewMockLoanBillingRepository(ctrl),
				loanRepo:    mock.NewMockLoanRepository(ctrl),
			}
			usecase := usecase.NewLoanUsecase(nil, mock.billingRepo, nil, mock.loanRepo, clock)
			if tc.mockFn != nil {
				tc.mockFn(mock)
			}

			err := usecase.ForceDisburse(context.Background(), tc.loanID)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
