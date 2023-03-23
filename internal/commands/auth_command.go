package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
)

type AuthCommand struct {
	usersService *services.UsersService
	credentials  *dto.Credentials
	context      *CommandContext
}

func NewAuthCommand(usersService *services.UsersService) (*AuthCommand, error) {
	credentials := &dto.Credentials{}
	context := &CommandContext{CommandContent: credentials}
	return &AuthCommand{usersService: usersService, context: context, credentials: credentials}, nil
}

func (a *AuthCommand) GetCommandContext() *CommandContext {
	return a.context
}

func (a *AuthCommand) Validate() error {
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

func (a *AuthCommand) Execute() (interface{}, error) {
	user, err := a.usersService.GetUserByUsername(a.credentials.Username)
	if err != nil {
		return nil, err
	}
	token := &dto.TokenResponse{
		Token: user.Id,
	}
	return token, nil
}