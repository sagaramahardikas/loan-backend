package config

import (
	"database/sql"
	"fmt"
	"net/http"

	"example.com/loan/module/user/internal/handler"
	"example.com/loan/module/user/internal/repository"
	"example.com/loan/module/user/internal/usecase"
)

type UserConfig struct {
	Database *sql.DB
}

func RegisterUserGatewayHandler(mux *http.ServeMux, cfg UserConfig) error {
	if cfg.Database == nil {
		return fmt.Errorf("database is not initialized")
	}

	userRepo := repository.NewUserRepository(cfg.Database)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)
	// Check Latest User Data (isDelinquent could be check in here through user status)
	mux.HandleFunc("/users/{id}", userHandler.GetUser())

	return nil
}
