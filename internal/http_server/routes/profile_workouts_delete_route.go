package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type ProfileWorkoutsDeleteRoute struct {
	userWorkoutsService *services.UserWorkoutsService
}

func NewProfileWorkoutsDeleteRoute(userWorkoutsService *services.UserWorkoutsService) *ProfileWorkoutsDeleteRoute {
	return &ProfileWorkoutsDeleteRoute{userWorkoutsService: userWorkoutsService}
}

func (*ProfileWorkoutsDeleteRoute) Method() string {
	return http.MethodDelete
}

func (*ProfileWorkoutsDeleteRoute) Pattern() string {
	return "/profile/workouts/active"
}

func (a *ProfileWorkoutsDeleteRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewProfileWorkoutsDeleteCommand(a.userWorkoutsService)
}
