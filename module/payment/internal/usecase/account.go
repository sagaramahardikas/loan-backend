package usecase

import (
	"context"

	"example.com/loan/module/payment/entity"
	"example.com/loan/module/payment/internal/repository"
)

type AccountUsecase interface {
	Create(ctx context.Context, account *entity.Account) error
}

type accountUsecase struct {
	repository repository.AccountRepository
}

func (u *accountUsecase) Create(ctx context.Context, account *entity.Account) error {
	return u.repository.Create(ctx, account)
}

func NewAccountUsecase(
	repository repository.AccountRepository,
) AccountUsecase {
	return &accountUsecase{
		repository: repository,
	}
}
