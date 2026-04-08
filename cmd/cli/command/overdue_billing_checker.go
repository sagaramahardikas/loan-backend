package command

import (
	"context"

	"example.com/loan/cmd/config"
	loanConfig "example.com/loan/module/loan/config"
)

type OverdueBillingCheckerCmd struct {
}

func (c *OverdueBillingCheckerCmd) Run() error {
	serviceConfig, err := config.LoadConfig()
	if err != nil {
		return err
	}

	db, err := config.InitializeDatabase(serviceConfig)
	if err != nil {
		return err
	}

	var cfg loanConfig.LoanConfig
	config.LoadLoanConfig(&cfg)
	cfg.Database = db
	dependencies, err := loanConfig.InitializeOverdueBillingCheckerDependencies(cfg)
	if err != nil {
		return err
	}

	return dependencies.Usecase.OverdueBillingChecker(context.Background())
}
