package client_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/loan/module/payment/client"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPaymentClient_CreateAndPayMutation(t *testing.T) {
	tests := []struct {
		name              string
		addressConfig     string
		useServerURL      bool
		mockServer        func() *httptest.Server
		expectedResponse  client.CreateAndPayMutationResponse
		expectedErrString string
	}{
		{
			name:          "error: error creating request",
			addressConfig: string([]byte{0x7f}),
			useServerURL:  false,
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			expectedErrString: "invalid control character in URL",
		},
		{
			name:          "error: error sending request",
			addressConfig: "htt://invalid-address",
			useServerURL:  false,
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			expectedErrString: "unsupported protocol scheme",
		},
		{
			name:         "error: non 200 status response",
			useServerURL: true,
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnprocessableEntity)
				}))
			},
			expectedErrString: "unexpected status code",
		},
		{
			name:         "error: decode body response",
			useServerURL: true,
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, `invalid json`)
				}))
			},
			expectedErrString: "invalid character 'i' looking for beginning of value",
		},
		{
			name:         "success",
			useServerURL: true,
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, `{"id":"1","account_id":"1","type":1,"reference":"REPAY-1","amount":"10000"}`)
				}))
			},
			expectedResponse: client.CreateAndPayMutationResponse{
				ID:        "1",
				AccountID: "1",
				Type:      1,
				Reference: "REPAY-1",
				Amount:    decimal.NewFromInt(10000),
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := tt.mockServer()
			defer srv.Close()

			if tt.useServerURL {
				tt.addressConfig = srv.URL
			}

			c := client.NewPaymentClient(http.DefaultClient, tt.addressConfig)
			response, err := c.CreateAndPayMutation(context.Background(), client.CreateAndPayMutationRequest{})
			if err != nil {
				assert.NotEmpty(t, tt.expectedErrString)
				assert.ErrorContains(t, err, tt.expectedErrString)
				return
			}

			assert.Equal(t, tt.expectedResponse, response)
		})
	}
}
