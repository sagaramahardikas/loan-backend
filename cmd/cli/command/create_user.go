package command

import (
	"context"

	"example.com/loan/cmd/config"
	usrConfig "example.com/loan/module/user/config"
	"example.com/loan/module/user/entity"
)

type CreateUserCmd struct {
	Username string `required:"" name:"username" help:"Username of the user" type:"string"`
	Status   string `required:"" name:"status" help:"Status of the user" type:"string"`
}

func (c *CreateUserCmd) Run() error {
	serviceConfig, err := config.LoadConfig()
	if err != nil {
		return err
	}

	db, err := config.InitializeDatabase(serviceConfig)
	if err != nil {
		return err
	}

	cfg := usrConfig.UserConfig{
		Database: db,
	}

	dependencies, err := usrConfig.InitializeCreateUserDependencies(cfg)
	if err != nil {
		return err
	}

	userStatus, err := entity.UserStatusString(c.Status)
	if err != nil {
		return err
	}

	return dependencies.Usecase.Create(context.Background(), entity.User{
		Username: c.Username,
		Status:   userStatus,
	})
}
