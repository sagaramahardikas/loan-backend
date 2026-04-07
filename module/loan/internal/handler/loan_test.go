package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/loan/module/loan/entity"
	"example.com/loan/module/loan/internal/handler"
	"example.com/loan/module/loan/internal/usecase/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockLoanHandler struct {
	usecase *mock.MockLoanUsecase
}

func TestLoanHandler_GetOutstandingLoans(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		mockFn         func(mock *mockLoanHandler)
		expectedResult string
	}{
		{
			name: "error: usecase error",
			id:   "123",
			mockFn: func(mocks *mockLoanHandler) {
				mocks.usecase.EXPECT().GetOutstandingLoans(
					gomock.Any(), "123",
				).Return(entity.GetOutstandingLoansResponse{}, errors.New("usecase error"))
			},
			expectedResult: "usecase error\n",
		},
		{
			name: "success: found",
			id:   "123",
			mockFn: func(mocks *mockLoanHandler) {
				mocks.usecase.EXPECT().GetOutstandingLoans(
					gomock.Any(), "123",
				).Return(entity.GetOutstandingLoansResponse{
					TotalOutstandingAmount: decimal.NewFromInt(100000),
					TotalBillingCount:      10,
				}, nil)
			},
			expectedResult: "{\"total_outstanding_amount\":\"100000\",\"total_billing_count\":10}\n",
		},
	}

	ctrl := gomock.NewController(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockLoanHandler{
				usecase: mock.NewMockLoanUsecase(ctrl),
			}

			r := httptest.NewRequest(http.MethodGet, "http://localhost/users/123", nil)
			r.SetPathValue("id", tc.id)
			w := httptest.NewRecorder()
			handler := handler.NewLoanHandler(mock.usecase)
			if tc.mockFn != nil {
				tc.mockFn(mock)
			}

			handler.GetOutstandingLoans()(w, r)
			body, _ := io.ReadAll(w.Body)
			responseText := string(body)
			assert.Equal(t, tc.expectedResult, responseText)
		})
	}
}

func TestLoanHandler_PayBilling(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		input          entity.PayBillingRequest
		mockFn         func(mock *mockLoanHandler)
		expectedResult string
	}{
		{
			name: "error: usecase error",
			id:   "123",
			input: entity.PayBillingRequest{
				Amount: decimal.NewFromInt(10000),
			},
			mockFn: func(mocks *mockLoanHandler) {
				mocks.usecase.EXPECT().PayBilling(
					gomock.Any(), entity.PayBillingRequest{
						UserID:    "1",
						BillingID: "123",
						Amount:    decimal.NewFromInt(10000),
					},
				).Return(entity.LoanBilling{}, errors.New("usecase error"))
			},
			expectedResult: "usecase error\n",
		},
		{
			name: "success: pay billing",
			id:   "123",
			input: entity.PayBillingRequest{
				Amount: decimal.NewFromInt(10000),
			},
			mockFn: func(mocks *mockLoanHandler) {
				mocks.usecase.EXPECT().PayBilling(
					gomock.Any(), entity.PayBillingRequest{
						UserID:    "1",
						BillingID: "123",
						Amount:    decimal.NewFromInt(10000),
					},
				).Return(entity.LoanBilling{
					ID:      "123",
					LoanID:  "123",
					Amount:  decimal.NewFromInt(10000),
					Status:  entity.LoanBillingStatusPaid,
					DueDate: time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			expectedResult: "{\"id\":\"123\",\"loan_id\":\"123\",\"amount\":\"10000\",\"status\":2,\"due_date\":\"2024-06-30T00:00:00Z\"}\n",
		},
	}

	ctrl := gomock.NewController(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockLoanHandler{
				usecase: mock.NewMockLoanUsecase(ctrl),
			}

			payload, _ := json.Marshal(tc.input)
			r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost/loans/billings/%s/pay", tc.id), bytes.NewBuffer(payload))
			r.SetPathValue("id", tc.id)
			w := httptest.NewRecorder()
			handler := handler.NewLoanHandler(mock.usecase)
			if tc.mockFn != nil {
				tc.mockFn(mock)
			}

			handler.PayBilling()(w, r)
			body, _ := io.ReadAll(w.Body)
			responseText := string(body)
			assert.Equal(t, tc.expectedResult, responseText)
		})
	}
}
