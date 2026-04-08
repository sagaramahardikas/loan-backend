package usecase

import (
	"context"
	"fmt"

	"example.com/loan/module/loan/entity"
	"example.com/loan/module/loan/internal/repository"
	"example.com/loan/module/user/client"
)

type LoanBillingUsecase interface {
	OverdueBillingChecker(ctx context.Context) error
}

type loanBillingUsecase struct {
	userService     client.UserService
	loanBillingRepo repository.LoanBillingRepository
}

func (u *loanBillingUsecase) OverdueBillingChecker(ctx context.Context) error {
	billings, err := u.loanBillingRepo.OverdueBillings(ctx)
	if err != nil {
		return err
	}

	var userOverdueBillings = make(map[string]int)
	for _, billing := range billings {
		userOverdueBillings[billing.UserID]++
		if billing.Status == entity.LoanBillingStatusCreated {
			billing.Status = entity.LoanBillingStatusOverdue
			if err := u.loanBillingRepo.Update(ctx, billing); err != nil {
				return err
			}
		}
	}

	for userID, count := range userOverdueBillings {
		if count >= 2 {
			if _, err := u.userService.Update(ctx, client.UpdateUserRequest{
				UserID: userID,
				Status: 3,
			}); err != nil {
				fmt.Println("error updating user status:", err)
				return err
			}
		}
	}

	return nil
}

func NewLoanBillingUsecase(
	loanBillingRepo repository.LoanBillingRepository,
	userService client.UserService,
) LoanBillingUsecase {
	return &loanBillingUsecase{
		loanBillingRepo: loanBillingRepo,
		userService:     userService,
	}
}
