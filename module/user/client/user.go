package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserService interface {
	Update(ctx context.Context, req UpdateUserRequest) (UpdateUserResponse, error)
}

type UserClient struct {
	Address string
	Client  *http.Client
}

func (c *UserClient) Update(ctx context.Context, req UpdateUserRequest) (UpdateUserResponse, error) {
	payload, _ := json.Marshal(req)
	request, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/internal/user/%s", c.Address, req.UserID), bytes.NewBuffer(payload))
	if err != nil {
		return UpdateUserResponse{}, err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(request)
	if err != nil {
		return UpdateUserResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UpdateUserResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response UpdateUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return UpdateUserResponse{}, err
	}

	return response, nil
}

func NewUserClient(client *http.Client, address string) *UserClient {
	return &UserClient{
		Address: address,
		Client:  client,
	}
}
