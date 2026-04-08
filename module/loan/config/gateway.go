package config

import (
	"database/sql"
	"fmt"
	"net/http"

	"example.com/loan/module/loan/internal/handler"
	"example.com/loan/module/loan/internal/repository"
	"example.com/loan/module/loan/internal/usecase"
	"example.com/loan/module/payment/client"
)

type LoanConfig struct {
	Database *sql.DB

	PaymentServiceAddress string `env:"PAYMENT_SERVICE_ADDRESS" default:"http://localhost:8888"`
	UserServiceAddress    string `env:"USER_SERVICE_ADDRESS" default:"http://localhost:8888"`
}

func RegisterLoanGatewayHandler(mux *http.ServeMux, cfg LoanConfig) error {
	if cfg.Database == nil {
		return fmt.Errorf("database is not initialized")
	}

	paymentClient := client.NewPaymentClient(http.DefaultClient, cfg.PaymentServiceAddress)
	loanBillingRepo := repository.NewLoanBillingRepository(cfg.Database)
	repaymentRepo := repository.NewRepaymentRepository(cfg.Database)
	loanRepo := repository.NewLoanRepository(cfg.Database)
	loanUsecase := usecase.NewLoanUsecase(paymentClient, loanBillingRepo, repaymentRepo, loanRepo, nil)
	loanHandler := handler.NewLoanHandler(loanUsecase)

	mux.HandleFunc("/loans/{id}/outstanding", loanHandler.GetOutstandingLoans())
	mux.HandleFunc("POST /loans/billings/{id}/pay", loanHandler.PayBilling())

	return nil
}
