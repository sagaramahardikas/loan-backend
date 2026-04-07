package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/loan/module/payment/entity"
	"example.com/loan/module/payment/internal/handler"
	"example.com/loan/module/payment/internal/usecase/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockMutationHandler struct {
	usecase *mock.MockMutationUsecase
}

func TestMutationHandler_Debit(t *testing.T) {
	testCases := []struct {
		name           string
		input          entity.Mutation
		mockFn         func(mock *mockMutationHandler)
		expectedResult string
	}{
		{
			name: "error: usecase error",
			input: entity.Mutation{
				AccountID: "123",
				Amount:    decimal.NewFromInt(10000),
				Type:      entity.MutationTypeRepayment,
				Reference: "ref-1",
			},
			mockFn: func(mocks *mockMutationHandler) {
				mocks.usecase.EXPECT().Debit(
					gomock.Any(), entity.Mutation{
						AccountID: "123",
						Amount:    decimal.NewFromInt(10000),
						Type:      entity.MutationTypeRepayment,
						Reference: "ref-1",
					},
				).Return(entity.Mutation{}, errors.New("usecase error"))
			},
			expectedResult: "usecase error\n",
		},
		{
			name: "success: create and pay mutation",
			input: entity.Mutation{
				AccountID: "123",
				Amount:    decimal.NewFromInt(10000),
				Type:      entity.MutationTypeRepayment,
				Reference: "ref-1",
			},
			mockFn: func(mocks *mockMutationHandler) {
				mocks.usecase.EXPECT().Debit(
					gomock.Any(), entity.Mutation{
						AccountID: "123",
						Amount:    decimal.NewFromInt(10000),
						Type:      entity.MutationTypeRepayment,
						Reference: "ref-1",
					},
				).Return(entity.Mutation{
					ID:        "1",
					AccountID: "123",
					Amount:    decimal.NewFromInt(10000),
					Type:      entity.MutationTypeRepayment,
					Reference: "ref-1",
				}, nil)
			},
			expectedResult: "{\"id\":\"1\",\"account_id\":\"123\",\"type\":1,\"reference\":\"ref-1\",\"amount\":\"10000\"}\n",
		},
	}

	ctrl := gomock.NewController(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockMutationHandler{
				usecase: mock.NewMockMutationUsecase(ctrl),
			}

			postBody := tc.input
			payload, _ := json.Marshal(postBody)

			r := httptest.NewRequest(http.MethodPost, "http://localhost/internal/payment/mutation-payments", bytes.NewBuffer(payload))
			w := httptest.NewRecorder()
			handler := handler.NewMutationHandler(mock.usecase)
			if tc.mockFn != nil {
				tc.mockFn(mock)
			}

			handler.CreateAndPayMutation()(w, r)
			body, _ := io.ReadAll(w.Body)
			responseText := string(body)
			assert.Equal(t, tc.expectedResult, responseText)
		})
	}
}
