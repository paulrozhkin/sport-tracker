package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type WorkoutsGetRoute struct {
	workoutsService *services.WorkoutsService
}

func NewWorkoutsGetRoute(workoutsService *services.WorkoutsService) *WorkoutsGetRoute {
	return &WorkoutsGetRoute{workoutsService: workoutsService}
}

func (*WorkoutsGetRoute) Method() string {
	return http.MethodGet
}

func (*WorkoutsGetRoute) Pattern() string {
	return "/workouts"
}

func (a *WorkoutsGetRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewWorkoutsGetCommandCommand(a.workoutsService)
}
