package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type ProfileWorkoutsCalendarGetByIdRoute struct {
	calendarService *services.UserWorkoutsCalendarService
}

func NewProfileWorkoutsCalendarGetByIdRoute(calendarService *services.UserWorkoutsCalendarService) *ProfileWorkoutsCalendarGetByIdRoute {
	return &ProfileWorkoutsCalendarGetByIdRoute{calendarService: calendarService}
}

func (*ProfileWorkoutsCalendarGetByIdRoute) Method() string {
	return http.MethodGet
}

func (*ProfileWorkoutsCalendarGetByIdRoute) Pattern() string {
	return "/profile/workouts/calendar/{workoutStatisticId}"
}

func (a *ProfileWorkoutsCalendarGetByIdRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewProfileWorkoutsCalendarGetByIdCommand(a.calendarService)
}
