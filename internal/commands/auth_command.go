package commands

import (
	"errors"
	"fmt"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type AuthCommand struct {
	UnauthorizedCommand
	usersService *services.UsersService
	credentials  *dto.Credentials
	context      *CommandContext
	tokenService *services.TokenService
}

func NewAuthCommand(usersService *services.UsersService, tokenService *services.TokenService) (*AuthCommand, error) {
	credentials := &dto.Credentials{}
	context := &CommandContext{CommandContent: credentials}
	return &AuthCommand{
		usersService: usersService,
		context:      context,
		credentials:  credentials,
		tokenService: tokenService}, nil
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
	if !utils.CheckPasswordHash(a.credentials.Password, user.Password) {
		return nil, models.NewNotFoundError("users", fmt.Sprintf("(%s, secret)", user.Username),
			"(username, password)")
	}
	token, err := a.tokenService.CreateToken(user)
	if err != nil {
		return nil, err
	}
	return &dto.TokenResponse{
		Token: token,
	}, nil
}
