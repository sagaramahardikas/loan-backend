package usecase

import (
	"context"

	"example.com/loan/module/loan/entity"
	"example.com/loan/module/loan/internal/repository"
)

type LoanUsecase interface {
	GetOutstandingLoans(ctx context.Context, loanID string) (entity.GetOutstandingLoansResponse, error)
}

type loanUsecase struct {
	loanBillingRepo repository.LoanBillingRepository
}

func (u *loanUsecase) GetOutstandingLoans(ctx context.Context, loanID string) (entity.GetOutstandingLoansResponse, error) {
	response, err := u.loanBillingRepo.SumOutstandingLoans(ctx, loanID)
	if err != nil {
		return entity.GetOutstandingLoansResponse{}, err
	}

	return response, nil
}

func NewLoanUsecase(loanBillingRepo repository.LoanBillingRepository) LoanUsecase {
	return &loanUsecase{
		loanBillingRepo: loanBillingRepo,
	}
}
