package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type ExercisesCreateRoute struct {
	exercisesService *services.ExercisesService
}

func NewExercisesCreateRoute(exercisesService *services.ExercisesService) *ExercisesCreateRoute {
	return &ExercisesCreateRoute{exercisesService: exercisesService}
}

func (*ExercisesCreateRoute) Method() string {
	return http.MethodPost
}

func (*ExercisesCreateRoute) Pattern() string {
	return "/exercises"
}

func (a *ExercisesCreateRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewExercisesCreateCommand(a.exercisesService)
}
