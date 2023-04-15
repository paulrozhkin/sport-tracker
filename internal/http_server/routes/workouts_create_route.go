package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type WorkoutsCreateRoute struct {
	workoutsService *services.WorkoutsService
}

func NewWorkoutsCreateRoute(workoutsService *services.WorkoutsService) *WorkoutsCreateRoute {
	return &WorkoutsCreateRoute{workoutsService: workoutsService}
}

func (*WorkoutsCreateRoute) Method() string {
	return http.MethodPost
}

func (*WorkoutsCreateRoute) Pattern() string {
	return "/workouts"
}

func (a *WorkoutsCreateRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewWorkoutsCreateCommand(a.workoutsService)
}
