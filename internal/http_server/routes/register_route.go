package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/metrics"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type RegisterRoute struct {
	userService  *services.UsersService
	tokenService *services.TokenService
	usersMetrics *metrics.UsersMetrics
}

func NewRegisterRoute(userService *services.UsersService,
	tokenService *services.TokenService,
	usersMetrics *metrics.UsersMetrics) *RegisterRoute {
	return &RegisterRoute{userService: userService,
		tokenService: tokenService,
		usersMetrics: usersMetrics}
}

func (*RegisterRoute) Method() string {
	return http.MethodPost
}

func (*RegisterRoute) Pattern() string {
	return "/register"
}

func (a *RegisterRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewRegisterCommand(a.userService, a.tokenService, a.usersMetrics)
}
