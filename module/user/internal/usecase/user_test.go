package usecase_test

import (
	"context"
	"errors"
	"testing"

	"example.com/loan/module/user/entity"
	"example.com/loan/module/user/internal/repository/mock"
	"example.com/loan/module/user/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockUserUsecase struct {
	repository *mock.MockUserRepository
}

func TestUserUsecase_GetByID(t *testing.T) {
	user := entity.User{
		ID:       "123",
		Username: "testuser",
		Status:   entity.UserStatusActive,
	}

	testCases := []struct {
		name         string
		id           string
		mockFn       func(mock *mockUserUsecase)
		expectedUser entity.User
		expectedErr  error
	}{
		{
			name: "error: db connection error",
			id:   "123",
			mockFn: func(mocks *mockUserUsecase) {
				mocks.repository.EXPECT().GetByID(
					gomock.Any(), "123",
				).Return(entity.User{}, errors.New("db connection error"))
			},
			expectedErr: errors.New("db connection error"),
		},
		{
			name: "success: found",
			id:   "123",
			mockFn: func(mocks *mockUserUsecase) {
				mocks.repository.EXPECT().GetByID(
					gomock.Any(), "123",
				).Return(user, nil)
			},
			expectedUser: user,
		},
	}

	ctrl := gomock.NewController(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockUserUsecase{
				repository: mock.NewMockUserRepository(ctrl),
			}

			usecase := usecase.NewUserUsecase(mock.repository)
			if tc.mockFn != nil {
				tc.mockFn(mock)
			}

			user, err := usecase.GetByID(context.Background(), tc.id)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestUserUsecase_Create(t *testing.T) {
	user := entity.User{
		ID:       "1",
		Username: "testuser",
		Status:   entity.UserStatusActive,
	}

	testCases := []struct {
		name        string
		input       entity.User
		mockFn      func(mock *mockUserUsecase)
		expectedErr error
	}{
		{
			name:  "error: db connection error",
			input: user,
			mockFn: func(mocks *mockUserUsecase) {
				mocks.repository.EXPECT().Create(
					gomock.Any(), &user,
				).Return(errors.New("db connection error"))
			},
			expectedErr: errors.New("db connection error"),
		},
		{
			name:  "success: created",
			input: user,
			mockFn: func(mocks *mockUserUsecase) {
				mocks.repository.EXPECT().Create(
					gomock.Any(), &user,
				).Return(nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockUserUsecase{repository: mock.NewMockUserRepository(ctrl)}
			usecase := usecase.NewUserUsecase(mock.repository)
			if tc.mockFn != nil {
				tc.mockFn(mock)
			}

			err := usecase.Create(context.Background(), tc.input)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
