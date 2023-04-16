package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type WorkoutPlansUpdateByIdRoute struct {
	workoutPlansService *services.WorkoutPlansService
}

func NewWorkoutPlansUpdateByIdRoute(workoutPlansService *services.WorkoutPlansService) *WorkoutPlansUpdateByIdRoute {
	return &WorkoutPlansUpdateByIdRoute{workoutPlansService: workoutPlansService}
}

func (*WorkoutPlansUpdateByIdRoute) Method() string {
	return http.MethodPut
}

func (*WorkoutPlansUpdateByIdRoute) Pattern() string {
	return "/workoutPlans/{workoutPlanId}"
}

func (a *WorkoutPlansUpdateByIdRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewWorkoutPlansUpdateByIdCommand(a.workoutPlansService)
}
