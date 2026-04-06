package config

import (
	"net/http"

	"example.com/loan/module/user/internal/handler"
	"example.com/loan/module/user/internal/repository"
	"example.com/loan/module/user/internal/usecase"
)

func RegisterUserGatewayHandler(mux *http.ServeMux) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	db, err := initializeDatabase(cfg)
	if err != nil {
		return err
	}

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)
	// Check Latest User Data (isDelinquent could be check in here through user status)
	mux.HandleFunc("/users/{id}", userHandler.GetUser())

	return nil
}
