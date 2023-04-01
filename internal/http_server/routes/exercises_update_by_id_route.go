package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type ExercisesUpdateByIdRoute struct {
	exercisesService *services.ExercisesService
}

func NewExercisesUpdateByIdRoute(exercisesService *services.ExercisesService) *ExercisesUpdateByIdRoute {
	return &ExercisesUpdateByIdRoute{exercisesService: exercisesService}
}

func (*ExercisesUpdateByIdRoute) Method() string {
	return http.MethodPut
}

func (*ExercisesUpdateByIdRoute) Pattern() string {
	return "/exercises/{exerciseId}"
}

func (a *ExercisesUpdateByIdRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewExercisesUpdateByIdCommand(a.exercisesService)
}
