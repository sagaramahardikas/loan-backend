package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/loan/module/user/entity"
	"example.com/loan/module/user/internal/handler"
	"example.com/loan/module/user/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockUserHandler struct {
	usecase *mock.MockUserUsecase
}

func TestUserHandler_GetUser(t *testing.T) {
	user := entity.User{
		ID:       "123",
		Username: "testuser",
		Status:   entity.UserStatusActive,
	}

	testCases := []struct {
		name           string
		id             string
		mockFn         func(mock *mockUserHandler)
		expectedResult string
	}{
		{
			name: "error: usecase error",
			id:   "123",
			mockFn: func(mocks *mockUserHandler) {
				mocks.usecase.EXPECT().GetByID(
					gomock.Any(), "123",
				).Return(entity.User{}, errors.New("usecase error"))
			},
			expectedResult: "usecase error\n",
		},
		{
			name: "success: found",
			id:   "123",
			mockFn: func(mocks *mockUserHandler) {
				mocks.usecase.EXPECT().GetByID(
					gomock.Any(), "123",
				).Return(user, nil)
			},
			expectedResult: "{\"user\":{\"id\":\"123\",\"username\":\"testuser\",\"status\":2}}\n",
		},
	}

	ctrl := gomock.NewController(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockUserHandler{
				usecase: mock.NewMockUserUsecase(ctrl),
			}

			r := httptest.NewRequest(http.MethodGet, "http://localhost/users/123", nil)
			r.SetPathValue("id", tc.id)
			w := httptest.NewRecorder()
			handler := handler.NewUserHandler(mock.usecase)
			if tc.mockFn != nil {
				tc.mockFn(mock)
			}

			handler.GetUser()(w, r)
			body, _ := io.ReadAll(w.Body)
			responseText := string(body)
			assert.Equal(t, tc.expectedResult, responseText)
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	user := entity.User{
		ID:     "123",
		Status: entity.UserStatusActive,
	}

	testCases := []struct {
		name           string
		input          entity.User
		mockFn         func(mock *mockUserHandler)
		expectedResult string
	}{
		{
			name:  "error: usecase error",
			input: user,
			mockFn: func(mocks *mockUserHandler) {
				mocks.usecase.EXPECT().Update(
					gomock.Any(), user,
				).Return(errors.New("usecase error"))
			},
			expectedResult: "usecase error\n",
		},
		{
			name:  "success: update user",
			input: user,
			mockFn: func(mocks *mockUserHandler) {
				mocks.usecase.EXPECT().Update(
					gomock.Any(), user,
				).Return(nil)
			},
			expectedResult: "{\"id\":\"123\",\"username\":\"\",\"status\":2}\n",
		},
	}

	ctrl := gomock.NewController(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockUserHandler{
				usecase: mock.NewMockUserUsecase(ctrl),
			}

			payload, _ := json.Marshal(tc.input)
			r := httptest.NewRequest(http.MethodPut, "http://localhost/users/123", bytes.NewBuffer(payload))
			r.SetPathValue("id", tc.input.ID)
			w := httptest.NewRecorder()
			handler := handler.NewUserHandler(mock.usecase)
			if tc.mockFn != nil {
				tc.mockFn(mock)
			}

			handler.UpdateUser()(w, r)
			body, _ := io.ReadAll(w.Body)
			responseText := string(body)
			assert.Equal(t, tc.expectedResult, responseText)
		})
	}
}
