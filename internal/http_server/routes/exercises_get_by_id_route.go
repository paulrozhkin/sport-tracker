package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type ExercisesGetByIdRoute struct {
	exercisesService *services.ExercisesService
}

func NewExercisesGetByIdRoute(exercisesService *services.ExercisesService) *ExercisesGetByIdRoute {
	return &ExercisesGetByIdRoute{exercisesService: exercisesService}
}

func (*ExercisesGetByIdRoute) Method() string {
	return http.MethodGet
}

func (*ExercisesGetByIdRoute) Pattern() string {
	return "/exercises/{exerciseId}"
}

func (a *ExercisesGetByIdRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewExercisesGetByIdCommand(a.exercisesService)
}
