package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"example.com/loan/module/payment/entity"
	"example.com/loan/module/payment/internal/usecase"
)

type MutationHandler struct {
	usecase usecase.MutationUsecase
}

func (h *MutationHandler) CreateAndPayMutation() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var debitReq entity.DebitRequest
		if err := json.NewDecoder(r.Body).Decode(&debitReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		debitReq.Type = entity.MutationTypeRepayment
		response, err := h.usecase.Debit(context.Background(), debitReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func NewMutationHandler(usecase usecase.MutationUsecase) *MutationHandler {
	return &MutationHandler{
		usecase: usecase,
	}
}
