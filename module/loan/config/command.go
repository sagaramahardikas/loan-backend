package config

import (
	"fmt"
	"net/http"

	"example.com/loan/module/loan/internal/repository"
	"example.com/loan/module/loan/internal/usecase"
	"example.com/loan/module/user/client"
)

type CreateLoanDependencies struct {
	Usecase usecase.LoanUsecase
}

func InitializeCreateLoanDependencies(cfg LoanConfig) (*CreateLoanDependencies, error) {
	if cfg.Database == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	loanRepo := repository.NewLoanRepository(cfg.Database)
	loanUsecase := usecase.NewLoanUsecase(nil, nil, nil, loanRepo, nil)

	return &CreateLoanDependencies{
		Usecase: loanUsecase,
	}, nil
}

type ForceDisburseLoanDependencies struct {
	Usecase usecase.LoanUsecase
}

func InitializeForceDisburseLoanDependencies(cfg LoanConfig) (*ForceDisburseLoanDependencies, error) {
	if cfg.Database == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	loanRepo := repository.NewLoanRepository(cfg.Database)
	loanBillingRepo := repository.NewLoanBillingRepository(cfg.Database)
	loanUsecase := usecase.NewLoanUsecase(nil, loanBillingRepo, nil, loanRepo, nil)

	return &ForceDisburseLoanDependencies{
		Usecase: loanUsecase,
	}, nil
}

type OverdueBillingCheckerDependencies struct {
	Usecase usecase.LoanBillingUsecase
}

func InitializeOverdueBillingCheckerDependencies(cfg LoanConfig) (*OverdueBillingCheckerDependencies, error) {
	if cfg.Database == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	loanBillingRepo := repository.NewLoanBillingRepository(cfg.Database)
	userService := client.NewUserClient(http.DefaultClient, cfg.UserServiceAddress)
	loanBillingUsecase := usecase.NewLoanBillingUsecase(loanBillingRepo, userService)

	return &OverdueBillingCheckerDependencies{
		Usecase: loanBillingUsecase,
	}, nil
}
