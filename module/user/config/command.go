package config

import (
	"fmt"

	"example.com/loan/module/user/internal/repository"
	"example.com/loan/module/user/internal/usecase"
)

type CreateUserDependencies struct {
	Usecase usecase.UserUsecase
}

func InitializeCreateUserDependencies(cfg UserConfig) (*CreateUserDependencies, error) {
	if cfg.Database == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	userRepo := repository.NewUserRepository(cfg.Database)
	userUsecase := usecase.NewUserUsecase(userRepo)

	return &CreateUserDependencies{
		Usecase: userUsecase,
	}, nil
}
