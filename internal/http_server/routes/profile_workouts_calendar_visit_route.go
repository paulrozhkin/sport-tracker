package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type ProfileWorkoutsCalendarVisitRoute struct {
	userWorkoutsCalendarService *services.UserWorkoutsCalendarService
}

func NewProfileWorkoutsCalendarVisitRoute(userWorkoutsCalendarService *services.UserWorkoutsCalendarService) *ProfileWorkoutsCalendarVisitRoute {
	return &ProfileWorkoutsCalendarVisitRoute{userWorkoutsCalendarService: userWorkoutsCalendarService}
}

func (*ProfileWorkoutsCalendarVisitRoute) Method() string {
	return http.MethodPost
}

func (*ProfileWorkoutsCalendarVisitRoute) Pattern() string {
	return "/profile/workouts/calendar/{workoutVisitId}"
}

func (a *ProfileWorkoutsCalendarVisitRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewProfileWorkoutsCalendarVisitCommand(a.userWorkoutsCalendarService)
}
