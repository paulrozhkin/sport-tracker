package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type WorkoutPlansGetRoute struct {
	workoutPlansService *services.WorkoutPlansService
}

func NewWorkoutPlansGetRoute(workoutPlansService *services.WorkoutPlansService) *WorkoutPlansGetRoute {
	return &WorkoutPlansGetRoute{workoutPlansService: workoutPlansService}
}

func (*WorkoutPlansGetRoute) Method() string {
	return http.MethodGet
}

func (*WorkoutPlansGetRoute) Pattern() string {
	return "/workoutPlans"
}

func (a *WorkoutPlansGetRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewWorkoutPlansGetCommandCommand(a.workoutPlansService)
}
