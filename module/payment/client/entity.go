package client

import "github.com/shopspring/decimal"

type CreateAndPayMutationRequest struct {
	UserID    string          `json:"user_id"`
	Amount    decimal.Decimal `json:"amount"`
	Reference string          `json:"reference"`
}

type CreateAndPayMutationResponse struct {
	ID        string          `json:"id"`
	AccountID string          `json:"account_id"`
	Type      int             `json:"type"`
	Reference string          `json:"reference"`
	Amount    decimal.Decimal `json:"amount"`
}
