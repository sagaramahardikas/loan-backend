package usecase

import (
	"context"
	"fmt"

	"example.com/loan/module/payment/entity"
	"example.com/loan/module/payment/internal/repository"
)

type MutationUsecase interface {
	Debit(ctx context.Context, mutation entity.Mutation) (entity.Mutation, error)
}

type mutationUsecase struct {
	accountRepo  repository.AccountRepository
	mutationRepo repository.MutationRepository
}

func (u *mutationUsecase) Debit(ctx context.Context, mutation entity.Mutation) (entity.Mutation, error) {
	// Improvement: put row lock in here, and wrap in transaction
	account, err := u.accountRepo.GetByUserID(ctx, mutation.AccountID)
	if err != nil {
		return entity.Mutation{}, err
	}

	if account.Status != entity.AccountStatusActive {
		return entity.Mutation{}, fmt.Errorf("account for user %s is %s", mutation.AccountID, account.Status)
	}

	// Check Mutation Exist or not
	existingMutation, err := u.mutationRepo.Get(ctx, mutation)
	if err != nil {
		return entity.Mutation{}, err
	} else if existingMutation.ID != "0" && existingMutation.ID != "" {
		return existingMutation, nil
	}

	if account.Balance.Sub(mutation.Amount).IsNegative() {
		return entity.Mutation{}, fmt.Errorf("insufficient balance for account %s", mutation.AccountID)
	}

	account.Balance = account.Balance.Sub(mutation.Amount)
	if err := u.accountRepo.Update(ctx, account); err != nil {
		return entity.Mutation{}, err
	}

	if err := u.mutationRepo.Create(ctx, &mutation); err != nil {
		return entity.Mutation{}, err
	}

	return mutation, nil
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
