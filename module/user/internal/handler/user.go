package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"example.com/loan/module/user/entity"
	"example.com/loan/module/user/internal/usecase"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func (h *UserHandler) GetUser() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		user, err := h.usecase.GetByID(context.Background(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := entity.GetResponse{User: user}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}
