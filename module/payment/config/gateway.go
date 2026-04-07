package config

import (
	"database/sql"
	"fmt"
	"net/http"

	"example.com/loan/module/payment/internal/handler"
	"example.com/loan/module/payment/internal/repository"
	"example.com/loan/module/payment/internal/usecase"
)

type PaymentConfig struct {
	Database *sql.DB
}

func RegisterPaymentGatewayHandler(mux *http.ServeMux, cfg PaymentConfig) error {
	if cfg.Database == nil {
		return fmt.Errorf("database is not initialized")
	}

	accountRepo := repository.NewAccountRepository(cfg.Database)
	mutationRepo := repository.NewMutationRepository(cfg.Database)
	mutationUsecase := usecase.NewMutationUsecase(accountRepo, mutationRepo)
	mutationHandler := handler.NewMutationHandler(mutationUsecase)

	mux.HandleFunc("POST /internal/payment/mutation-payments", mutationHandler.CreateAndPayMutation())

	return nil
}
