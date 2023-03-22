package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type AuthRoute struct {
	userService *services.UsersService
}

func NewAuthRoute(userService *services.UsersService) *AuthRoute {
	return &AuthRoute{userService: userService}
}

func (*AuthRoute) Method() string {
	return http.MethodPost
}

func (*AuthRoute) Pattern() string {
	return "/auth"
}

func (a *AuthRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewAuthCommand(a.userService)
}
