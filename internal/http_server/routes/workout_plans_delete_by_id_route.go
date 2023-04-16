package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type WorkoutPlansDeleteByIdRoute struct {
	workoutPlansService *services.WorkoutPlansService
}

func NewWorkoutPlansDeleteByIdRoute(workoutPlansService *services.WorkoutPlansService) *WorkoutPlansDeleteByIdRoute {
	return &WorkoutPlansDeleteByIdRoute{workoutPlansService: workoutPlansService}
}

func (*WorkoutPlansDeleteByIdRoute) Method() string {
	return http.MethodDelete
}

func (*WorkoutPlansDeleteByIdRoute) Pattern() string {
	return "/workoutPlans/{workoutPlanId}"
}

func (a *WorkoutPlansDeleteByIdRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewWorkoutPlansDeleteByIdCommand(a.workoutPlansService)
}
