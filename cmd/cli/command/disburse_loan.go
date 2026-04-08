package command

import (
	"context"

	"example.com/loan/cmd/config"
	loanConfig "example.com/loan/module/loan/config"
)

type ForceDisburseLoanCmd struct {
	LoanID string `required:"" name:"loan_id" help:"ID of the loan" type:"string"`
}

func (c *ForceDisburseLoanCmd) Run() error {
	serviceConfig, err := config.LoadConfig()
	if err != nil {
		return err
	}

	db, err := config.InitializeDatabase(serviceConfig)
	if err != nil {
		return err
	}

	cfg := loanConfig.LoanConfig{
		Database: db,
	}

	dependencies, err := loanConfig.InitializeForceDisburseLoanDependencies(cfg)
	if err != nil {
		return err
	}

	return dependencies.Usecase.ForceDisburse(context.Background(), c.LoanID)
}
