package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type RegisterRoute struct {
	userService *services.UsersService
}

func NewRegisterRoute(userService *services.UsersService) *RegisterRoute {
	return &RegisterRoute{userService: userService}
}

func (*RegisterRoute) Method() string {
	return http.MethodPost
}

func (*RegisterRoute) Pattern() string {
	return "/register"
}

func (a *RegisterRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewRegisterCommand(a.userService)
}
