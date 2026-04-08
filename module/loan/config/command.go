package config

import (
	"fmt"

	"example.com/loan/module/loan/internal/repository"
	"example.com/loan/module/loan/internal/usecase"
)

type CreateLoanDependencies struct {
	Usecase usecase.LoanUsecase
}

func InitializeCreateLoanDependencies(cfg LoanConfig) (*CreateLoanDependencies, error) {
	if cfg.Database == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	loanRepo := repository.NewLoanRepository(cfg.Database)
	loanUsecase := usecase.NewLoanUsecase(nil, nil, nil, loanRepo)

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
	loanUsecase := usecase.NewLoanUsecase(nil, loanBillingRepo, nil, loanRepo)

	return &ForceDisburseLoanDependencies{
		Usecase: loanUsecase,
	}, nil
}
