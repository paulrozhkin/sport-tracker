package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type ExercisesGetRoute struct {
	exercisesService *services.ExercisesService
}

func NewExercisesGetRoute(exercisesService *services.ExercisesService) *ExercisesGetRoute {
	return &ExercisesGetRoute{exercisesService: exercisesService}
}

func (*ExercisesGetRoute) Method() string {
	return http.MethodGet
}

func (*ExercisesGetRoute) Pattern() string {
	return "/exercises"
}

func (a *ExercisesGetRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewExercisesGetCommandCommand(a.exercisesService)
}
