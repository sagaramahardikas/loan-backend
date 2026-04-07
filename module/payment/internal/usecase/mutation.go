package usecase

import (
	"context"
	"fmt"

	"example.com/loan/module/payment/entity"
	"example.com/loan/module/payment/internal/repository"
)

type MutationUsecase interface {
	Debit(ctx context.Context, req entity.DebitRequest) (entity.Mutation, error)
}

type mutationUsecase struct {
	accountRepo  repository.AccountRepository
	mutationRepo repository.MutationRepository
}

func (u *mutationUsecase) Debit(ctx context.Context, req entity.DebitRequest) (entity.Mutation, error) {
	// Improvement: put row lock in here, and wrap in transaction
	account, err := u.accountRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return entity.Mutation{}, err
	}

	if account.Status != entity.AccountStatusActive {
		return entity.Mutation{}, fmt.Errorf("account for user %s is %s", account.ID, account.Status)
	}

	// Check Mutation Exist or not
	reqMutation := entity.Mutation{
		AccountID: account.ID,
		Type:      req.Type,
		Reference: req.Reference,
		Amount:    req.Amount,
	}

	existingMutation, err := u.mutationRepo.Get(ctx, reqMutation)
	if err != nil {
		return entity.Mutation{}, err
	} else if existingMutation.ID != "0" && existingMutation.ID != "" {
		return existingMutation, nil
	}

	if account.Balance.Sub(req.Amount).IsNegative() {
		return entity.Mutation{}, fmt.Errorf("insufficient balance for account %s", account.ID)
	}

	account.Balance = account.Balance.Sub(req.Amount)
	if err := u.accountRepo.Update(ctx, account); err != nil {
		return entity.Mutation{}, err
	}

	if err := u.mutationRepo.Create(ctx, &reqMutation); err != nil {
		return entity.Mutation{}, err
	}

	return reqMutation, nil
}

func NewMutationUsecase(
	accountRepo repository.AccountRepository,
	mutationRepo repository.MutationRepository,
) MutationUsecase {
	return &mutationUsecase{
		accountRepo:  accountRepo,
		mutationRepo: mutationRepo,
	}
}
