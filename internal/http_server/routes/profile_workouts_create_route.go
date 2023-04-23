package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type ProfileWorkoutsCreateRoute struct {
	userWorkoutsService *services.UserWorkoutsService
}

func NewProfileWorkoutsCreateRoute(userWorkoutsService *services.UserWorkoutsService) *ProfileWorkoutsCreateRoute {
	return &ProfileWorkoutsCreateRoute{userWorkoutsService: userWorkoutsService}
}

func (*ProfileWorkoutsCreateRoute) Method() string {
	return http.MethodPost
}

func (*ProfileWorkoutsCreateRoute) Pattern() string {
	return "/profile/workouts"
}

func (a *ProfileWorkoutsCreateRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewProfileWorkoutsCreateCommand(a.userWorkoutsService)
}
