package commands

import "github.com/paulrozhkin/sport-tracker/internal/models"

type UnauthorizedCommand struct {
}

func (*UnauthorizedCommand) RequireAuthorization() bool {
	return false
}

func (a *UnauthorizedCommand) SetAuthorization(_ *models.Claims) {
}
