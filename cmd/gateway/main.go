package main

import (
	"fmt"
	"net/http"
	"os"

	loanConfig "example.com/loan/module/loan/config"
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
	serviceConfig, err := loadConfig()
	if err != nil {
		return err
	}

	db, err := initializeDatabase(serviceConfig)
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

	loanCfg := loanConfig.LoanConfig{
		Database: db,
	}
	err = loanConfig.RegisterLoanGatewayHandler(mux, loanCfg)
	if err != nil {
		return err
	}

	return nil
}
