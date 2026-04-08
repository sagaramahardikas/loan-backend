package config

import (
	"fmt"

	"example.com/loan/module/payment/internal/repository"
	"example.com/loan/module/payment/internal/usecase"
)

type CreateAccountDependencies struct {
	Usecase usecase.AccountUsecase
}

func InitializeCreateAccountDependencies(cfg PaymentConfig) (*CreateAccountDependencies, error) {
	if cfg.Database == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	accountRepo := repository.NewAccountRepository(cfg.Database)
	accountUsecase := usecase.NewAccountUsecase(accountRepo)

	return &CreateAccountDependencies{
		Usecase: accountUsecase,
	}, nil
}
