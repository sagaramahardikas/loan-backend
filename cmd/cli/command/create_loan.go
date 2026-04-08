package command

import (
	"context"

	"example.com/loan/cmd/config"
	loanConfig "example.com/loan/module/loan/config"
	"example.com/loan/module/loan/entity"
	"github.com/shopspring/decimal"
)

type CreateLoanCmd struct {
	UserID    string  `required:"" name:"user_id" help:"ID of the user" type:"string"`
	Principal int     `required:"" name:"principal" help:"Principal amount of the loan" type:"int"`
	Term      int     `required:"" name:"term" help:"Term of the loan in weeks" type:"int"`
	Interest  float64 `required:"" name:"interest" help:"Interest rate of the loan" type:"float"`
}

func (c *CreateLoanCmd) Run() error {
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

	dependencies, err := loanConfig.InitializeCreateLoanDependencies(cfg)
	if err != nil {
		return err
	}

	principalDecimal := decimal.NewFromInt(int64(c.Principal))
	totalAmountDecimal := principalDecimal.Mul(decimal.NewFromFloat(1 + c.Interest))
	weeklyInstallmentDecimal := totalAmountDecimal.Div(decimal.NewFromInt(int64(c.Term)))
	return dependencies.Usecase.Create(context.Background(), &entity.Loan{
		UserID:            c.UserID,
		Principal:         principalDecimal,
		Term:              c.Term,
		Interest:          c.Interest,
		TotalAmount:       totalAmountDecimal,
		WeeklyInstallment: weeklyInstallmentDecimal,
		Status:            entity.LoanStatusProposed,
	})
}
