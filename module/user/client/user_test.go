package client_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/loan/module/user/client"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserClient_Update(t *testing.T) {
	tests := []struct {
		name              string
		addressConfig     string
		useServerURL      bool
		request           client.UpdateUserRequest
		mockServer        func() *httptest.Server
		expectedResponse  client.UpdateUserResponse
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
			request: client.UpdateUserRequest{
				UserID: "1",
				Status: 3,
			},
			mockServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, `{"id":"1","username":"","status":3}`)
				}))
			},
			expectedResponse: client.UpdateUserResponse{
				UserID: "1",
				Status: 3,
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

			c := client.NewUserClient(http.DefaultClient, tt.addressConfig)
			response, err := c.Update(context.Background(), tt.request)
			if err != nil {
				assert.NotEmpty(t, tt.expectedErrString)
				assert.ErrorContains(t, err, tt.expectedErrString)
				return
			}

			assert.Equal(t, tt.expectedResponse, response)
		})
	}
}
