package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"example.com/loan/module/loan/entity"
	"example.com/loan/module/loan/internal/usecase"
)

type LoanHandler struct {
	usecase usecase.LoanUsecase
}

func (h *LoanHandler) GetOutstandingLoans() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		response, err := h.usecase.GetOutstandingLoans(context.Background(), id)
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

func (h *LoanHandler) PayBilling() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req entity.PayBillingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := r.PathValue("id")
		req.UserID = "1" // hardcode for now, should get from jwt token
		req.BillingID = id
		response, err := h.usecase.PayBilling(context.Background(), req)
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

func NewLoanHandler(usecase usecase.LoanUsecase) *LoanHandler {
	return &LoanHandler{
		usecase: usecase,
	}
}
