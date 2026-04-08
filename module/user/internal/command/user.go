package command

import (
	"context"

	"example.com/loan/module/user/entity"
	"example.com/loan/module/user/internal/usecase"
)

// UserCommand is a collection of commands for user resource that will be used by CLI.
type UserCommand struct {
	userUsecase usecase.UserUsecase
}

func (c *UserCommand) CreateUser(ctx context.Context, username, status string) error {
	userStatus, err := entity.UserStatusString(status)
	if err != nil {
		return err
	}

	user := entity.User{
		Username: username,
		Status:   userStatus,
	}

	return c.userUsecase.Create(ctx, user)
}

func NewUserCommand(usecase usecase.UserUsecase) *UserCommand {
	return &UserCommand{
		userUsecase: usecase,
	}
}
