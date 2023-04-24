package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type ProfileWorkoutsCalendarGetRoute struct {
	userWorkoutsCalendarService *services.UserWorkoutsCalendarService
}

func NewProfileWorkoutsCalendarGetRoute(userWorkoutsCalendarService *services.UserWorkoutsCalendarService) *ProfileWorkoutsCalendarGetRoute {
	return &ProfileWorkoutsCalendarGetRoute{userWorkoutsCalendarService: userWorkoutsCalendarService}
}

func (*ProfileWorkoutsCalendarGetRoute) Method() string {
	return http.MethodGet
}

func (*ProfileWorkoutsCalendarGetRoute) Pattern() string {
	return "/profile/workouts/calendar"
}

func (a *ProfileWorkoutsCalendarGetRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewProfileWorkoutsCalendarGetCommand(a.userWorkoutsCalendarService)
}
