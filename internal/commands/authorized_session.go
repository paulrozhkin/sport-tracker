package commands

import "github.com/paulrozhkin/sport-tracker/internal/models"

type AuthorizedSession struct {
	*models.Claims
}

func (*AuthorizedSession) RequireAuthorization() bool {
	return true
}
