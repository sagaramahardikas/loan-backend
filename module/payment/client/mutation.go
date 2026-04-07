package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type PaymentService interface {
	CreateAndPayMutation(ctx context.Context, req CreateAndPayMutationRequest) (CreateAndPayMutationResponse, error)
}

type PaymentClient struct {
	Address string
	Client  *http.Client
}

func (c *PaymentClient) CreateAndPayMutation(ctx context.Context, req CreateAndPayMutationRequest) (CreateAndPayMutationResponse, error) {
	payload, _ := json.Marshal(req)
	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/internal/payment/mutation-payments", c.Address), bytes.NewBuffer(payload))
	if err != nil {
		return CreateAndPayMutationResponse{}, err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(request)
	if err != nil {
		return CreateAndPayMutationResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CreateAndPayMutationResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response CreateAndPayMutationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return CreateAndPayMutationResponse{}, err
	}

	return response, nil
}

func NewPaymentClient(client *http.Client, address string) *PaymentClient {
	return &PaymentClient{
		Address: address,
		Client:  client,
	}
}
