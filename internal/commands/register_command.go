package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
)

type RegisterCommand struct {
	usersService *services.UsersService
	credentials  *dto.Credentials
	context      *CommandContext
}

func NewRegisterCommand(usersService *services.UsersService) (*RegisterCommand, error) {
	credentials := &dto.Credentials{}
	context := &CommandContext{CommandContent: credentials}
	return &RegisterCommand{usersService: usersService, context: context, credentials: credentials}, nil
}

func (a *RegisterCommand) GetCommandContext() *CommandContext {
	return a.context
}

func (a *RegisterCommand) Validate() error {
	validationError := new(models.ValidationError)
	if a.credentials.Username == "" {
		validationError.AddError("username", errors.New(models.ArgumentNullOrEmptyError))
	}
	if a.credentials.Password == "" {
		validationError.AddError("password", errors.New(models.ArgumentNullOrEmptyError))
	}
	if validationError.HasErrors() {
		return validationError
	}
	return nil
}

func (a *RegisterCommand) Execute() (interface{}, error) {
	user := models.User{Username: a.credentials.Username, Password: a.credentials.Password}
	newUser, err := a.usersService.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return &dto.TokenResponse{Token: newUser.Id}, nil
}
