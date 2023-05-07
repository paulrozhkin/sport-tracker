package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
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

func (c *RegisterCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *RegisterCommand) Validate() error {
	validationError := new(models.ValidationError)
	if c.credentials.Username == "" {
		validationError.AddError("username", errors.New(models.ArgumentNullOrEmptyError))
	}
	if c.credentials.Password == "" {
		validationError.AddError("password", errors.New(models.ArgumentNullOrEmptyError))
	}
	if validationError.HasErrors() {
		return validationError
	}
	return nil
}

func (c *RegisterCommand) Execute() (interface{}, error) {
	user := models.User{Username: c.credentials.Username, Password: c.credentials.Password}
	newUser, err := c.usersService.CreateUser(user)
	if err != nil {
		return nil, err
	}
	token, err := c.tokenService.CreateToken(newUser)
	if err != nil {
		return nil, err
	}
	return &dto.TokenResponse{
		Token: token,
	}, nil
}
