package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type WorkoutsGetByIdRoute struct {
	workoutsService *services.WorkoutsService
}

func NewWorkoutsGetByIdRoute(workoutsService *services.WorkoutsService) *WorkoutsGetByIdRoute {
	return &WorkoutsGetByIdRoute{workoutsService: workoutsService}
}

func (*WorkoutsGetByIdRoute) Method() string {
	return http.MethodGet
}

func (*WorkoutsGetByIdRoute) Pattern() string {
	return "/workouts/{workoutId}"
}

func (a *WorkoutsGetByIdRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewWorkoutsGetByIdCommand(a.workoutsService)
}
