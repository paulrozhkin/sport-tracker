package commands

import "github.com/paulrozhkin/sport-tracker/internal/models"

type CommandContext struct {
	CommandContent    interface{}
	CommandParameters map[string]string
}

type ICommand interface {
	GetCommandContext() *CommandContext
	Validate() error
	Execute() (interface{}, error)
	RequireAuthorization() bool
	SetAuthorization(claims *models.Claims)
}
