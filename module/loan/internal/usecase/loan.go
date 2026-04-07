package usecase

import (
	"context"
	"fmt"

	"example.com/loan/module/loan/entity"
	"example.com/loan/module/loan/internal/repository"
	"example.com/loan/module/payment/client"
)

var (
	RepaymentReferencePrefix = "REPAY"
)

type LoanUsecase interface {
	GetOutstandingLoans(ctx context.Context, loanID string) (entity.GetOutstandingLoansResponse, error)
	PayBilling(ctx context.Context, req entity.PayBillingRequest) (entity.LoanBilling, error)
}

type loanUsecase struct {
	paymentService  client.PaymentService
	loanBillingRepo repository.LoanBillingRepository
	repaymentRepo   repository.RepaymentRepository
}

func (u *loanUsecase) GetOutstandingLoans(ctx context.Context, loanID string) (entity.GetOutstandingLoansResponse, error) {
	response, err := u.loanBillingRepo.SumOutstandingLoans(ctx, loanID)
	if err != nil {
		return entity.GetOutstandingLoansResponse{}, err
	}

	return response, nil
}

func (u *loanUsecase) PayBilling(ctx context.Context, req entity.PayBillingRequest) (entity.LoanBilling, error) {
	billing, err := u.loanBillingRepo.GetByID(ctx, req.BillingID)
	if err != nil {
		return entity.LoanBilling{}, err
	}

	if billing.Status == entity.LoanBillingStatusPaid {
		return billing, nil
	}

	repayment := entity.Repayment{
		LoanBillingID: billing.ID,
		Amount:        req.Amount,
		Status:        entity.RepaymentStatusCreated,
	}
	if err := u.repaymentRepo.Create(ctx, &repayment); err != nil {
		return entity.LoanBilling{}, err
	}

	_, err = u.paymentService.CreateAndPayMutation(ctx, client.CreateAndPayMutationRequest{
		UserID:    req.UserID,
		Amount:    req.Amount,
		Reference: fmt.Sprintf("%s-%s", RepaymentReferencePrefix, repayment.ID),
	})
	if err != nil {
		return entity.LoanBilling{}, err
	}

	billing.Status = entity.LoanBillingStatusPaid
	if err := u.loanBillingRepo.Update(ctx, billing); err != nil {
		return entity.LoanBilling{}, err
	}

	repayment.Reference = fmt.Sprintf("%s-%s", RepaymentReferencePrefix, repayment.ID)
	if err := u.repaymentRepo.Update(ctx, repayment); err != nil {
		return entity.LoanBilling{}, err
	}

	return billing, nil
}

func NewLoanUsecase(
	paymentService client.PaymentService,
	loanBillingRepo repository.LoanBillingRepository,
	repaymentRepo repository.RepaymentRepository,
) LoanUsecase {
	return &loanUsecase{
		paymentService:  paymentService,
		loanBillingRepo: loanBillingRepo,
		repaymentRepo:   repaymentRepo,
	}
}
