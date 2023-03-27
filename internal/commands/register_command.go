package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type RegisterCommand struct {
	UnauthorizedCommand
	usersService *services.UsersService
	credentials  *dto.Credentials
	context      *CommandContext
	tokenService *services.TokenService
}

func NewRegisterCommand(usersService *services.UsersService, tokenService *services.TokenService) (*RegisterCommand, error) {
	credentials := &dto.Credentials{}
	context := &CommandContext{CommandContent: credentials}
	return &RegisterCommand{usersService: usersService, context: context, credentials: credentials, tokenService: tokenService}, nil
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
	hashedPassword, err := utils.HashPassword(a.credentials.Password)
	if err != nil {
		return nil, err
	}
	user := models.User{Username: a.credentials.Username, Password: hashedPassword}
	newUser, err := a.usersService.CreateUser(user)
	if err != nil {
		return nil, err
	}
	token, err := a.tokenService.CreateToken(newUser)
	if err != nil {
		return nil, err
	}
	return &dto.TokenResponse{
		Token: token,
	}, nil
}
