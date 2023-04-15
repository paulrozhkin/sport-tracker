package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type WorkoutsDeleteByIdRoute struct {
	workoutsService *services.WorkoutsService
}

func NewWorkoutsDeleteByIdRoute(workoutsService *services.WorkoutsService) *WorkoutsDeleteByIdRoute {
	return &WorkoutsDeleteByIdRoute{workoutsService: workoutsService}
}

func (*WorkoutsDeleteByIdRoute) Method() string {
	return http.MethodDelete
}

func (*WorkoutsDeleteByIdRoute) Pattern() string {
	return "/workouts/{workoutId}"
}

func (a *WorkoutsDeleteByIdRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewWorkoutsDeleteByIdCommand(a.workoutsService)
}
