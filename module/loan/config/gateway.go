package config

import (
	"database/sql"
	"fmt"
	"net/http"

	"example.com/loan/module/loan/internal/handler"
	"example.com/loan/module/loan/internal/repository"
	"example.com/loan/module/loan/internal/usecase"
)

type LoanConfig struct {
	Database *sql.DB
}

func RegisterLoanGatewayHandler(mux *http.ServeMux, cfg LoanConfig) error {
	if cfg.Database == nil {
		return fmt.Errorf("database is not initialized")
	}

	loanBillingRepo := repository.NewLoanBillingRepository(cfg.Database)
	loanUsecase := usecase.NewLoanUsecase(loanBillingRepo)
	loanHandler := handler.NewLoanHandler(loanUsecase)

	mux.HandleFunc("/loans/{id}/outstanding", loanHandler.GetOutstandingLoans())

	return nil
}
