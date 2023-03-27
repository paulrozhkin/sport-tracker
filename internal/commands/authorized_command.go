package commands

import "github.com/paulrozhkin/sport-tracker/internal/models"

type AuthorizedCommand struct {
	claims *models.Claims
}

func (*AuthorizedCommand) RequireAuthorization() bool {
	return true
}

func (a *AuthorizedCommand) SetAuthorization(claims *models.Claims) {
	a.claims = claims
}
