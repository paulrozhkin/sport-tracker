package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type WorkoutsUpdateByIdRoute struct {
	workoutsService *services.WorkoutsService
}

func NewWorkoutsUpdateByIdRoute(workoutsService *services.WorkoutsService) *WorkoutsUpdateByIdRoute {
	return &WorkoutsUpdateByIdRoute{workoutsService: workoutsService}
}

func (*WorkoutsUpdateByIdRoute) Method() string {
	return http.MethodPut
}

func (*WorkoutsUpdateByIdRoute) Pattern() string {
	return "/workouts/{workoutId}"
}

func (a *WorkoutsUpdateByIdRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewWorkoutsUpdateByIdCommand(a.workoutsService)
}
