package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type WorkoutPlansGetByIdRoute struct {
	workoutPlansService *services.WorkoutPlansService
}

func NewWorkoutPlansGetByIdRoute(workoutPlansService *services.WorkoutPlansService) *WorkoutPlansGetByIdRoute {
	return &WorkoutPlansGetByIdRoute{workoutPlansService: workoutPlansService}
}

func (*WorkoutPlansGetByIdRoute) Method() string {
	return http.MethodGet
}

func (*WorkoutPlansGetByIdRoute) Pattern() string {
	return "/workoutPlans/{workoutPlanId}"
}

func (a *WorkoutPlansGetByIdRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewWorkoutPlansGetByIdCommand(a.workoutPlansService)
}
