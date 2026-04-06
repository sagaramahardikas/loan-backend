package usecase

import (
	"context"

	"example.com/loan/module/user/entity"
	"example.com/loan/module/user/internal/repository"
)

type UserUsecase interface {
	GetByID(ctx context.Context, id string) (entity.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func (u *userUsecase) GetByID(ctx context.Context, id string) (entity.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}
