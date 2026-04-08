package main

import (
	"fmt"
	"net/http"
	"os"

	"example.com/loan/cmd/config"
	loanConfig "example.com/loan/module/loan/config"
	paymentConfig "example.com/loan/module/payment/config"
	usrConfig "example.com/loan/module/user/config"
)

func main() {
	mux := http.NewServeMux()
	err := RegisterGatewayHandler(mux)
	if err != nil {
		fmt.Printf("failed to register handler: %v\n", err)
		return
	}

	err = http.ListenAndServe("0.0.0.0:8888", mux)
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
		os.Exit(1)
	}
}

func RegisterGatewayHandler(mux *http.ServeMux) error {
	serviceConfig, err := config.LoadConfig()
	if err != nil {
		return err
	}

	db, err := config.InitializeDatabase(serviceConfig)
	if err != nil {
		return err
	}

	userCfg := usrConfig.UserConfig{
		Database: db,
	}
	err = usrConfig.RegisterUserGatewayHandler(mux, userCfg)
	if err != nil {
		return err
	}

	loanCfg := loanConfig.LoanConfig{Database: db}

	err = config.LoadLoanConfig(&loanCfg)
	if err != nil {
		return err
	}

	err = loanConfig.RegisterLoanGatewayHandler(mux, loanCfg)
	if err != nil {
		return err
	}

	paymentCfg := paymentConfig.PaymentConfig{
		Database: db,
	}
	err = paymentConfig.RegisterPaymentGatewayHandler(mux, paymentCfg)
	if err != nil {
		return err
	}

	return nil
}
