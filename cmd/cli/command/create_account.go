package command

import (
	"context"

	"example.com/loan/cmd/config"
	paymentConfig "example.com/loan/module/payment/config"
	"example.com/loan/module/payment/entity"
	"github.com/shopspring/decimal"
)

type CreateAccountCmd struct {
	UserID  string `required:"" name:"user_id" help:"ID of the user" type:"string"`
	Balance int    `required:"" name:"balance" help:"Initial balance of the account" type:"int"`
	Status  string `required:"" name:"status" help:"Status of the account" type:"string"`
}

func (c *CreateAccountCmd) Run() error {
	serviceConfig, err := config.LoadConfig()
	if err != nil {
		return err
	}

	db, err := config.InitializeDatabase(serviceConfig)
	if err != nil {
		return err
	}

	cfg := paymentConfig.PaymentConfig{
		Database: db,
	}

	dependencies, err := paymentConfig.InitializeCreateAccountDependencies(cfg)
	if err != nil {
		return err
	}

	accountStatus, err := entity.AccountStatusString(c.Status)
	if err != nil {
		return err
	}

	balanceDecimal := decimal.NewFromInt(int64(c.Balance))
	return dependencies.Usecase.Create(context.Background(), &entity.Account{
		UserID:  c.UserID,
		Balance: balanceDecimal,
		Status:  accountStatus,
	})
}
