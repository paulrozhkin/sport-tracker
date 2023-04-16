package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type WorkoutPlansCreateRoute struct {
	workoutPlansService *services.WorkoutPlansService
}

func NewWorkoutPlansCreateRoute(workoutPlansService *services.WorkoutPlansService) *WorkoutPlansCreateRoute {
	return &WorkoutPlansCreateRoute{workoutPlansService: workoutPlansService}
}

func (*WorkoutPlansCreateRoute) Method() string {
	return http.MethodPost
}

func (*WorkoutPlansCreateRoute) Pattern() string {
	return "/workoutPlans"
}

func (a *WorkoutPlansCreateRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewWorkoutPlansCreateCommand(a.workoutPlansService)
}
