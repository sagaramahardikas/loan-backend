package usecase

import (
	"context"

	"example.com/loan/module/user/entity"
	"example.com/loan/module/user/internal/repository"
)

type UserUsecase interface {
	GetByID(ctx context.Context, id string) (entity.User, error)
	Create(ctx context.Context, user entity.User) error
	Update(ctx context.Context, user entity.User) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func (u *userUsecase) GetByID(ctx context.Context, id string) (entity.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func (u *userUsecase) Create(ctx context.Context, user entity.User) error {
	return u.userRepo.Create(ctx, &user)
}

func (u *userUsecase) Update(ctx context.Context, user entity.User) error {
	return u.userRepo.Update(ctx, user)
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}
