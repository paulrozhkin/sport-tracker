package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type ExercisesDeleteByIdRoute struct {
	exercisesService *services.ExercisesService
}

func NewExercisesDeleteByIdRoute(exercisesService *services.ExercisesService) *ExercisesDeleteByIdRoute {
	return &ExercisesDeleteByIdRoute{exercisesService: exercisesService}
}

func (*ExercisesDeleteByIdRoute) Method() string {
	return http.MethodDelete
}

func (*ExercisesDeleteByIdRoute) Pattern() string {
	return "/exercises/{exerciseId}"
}

func (a *ExercisesDeleteByIdRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewExercisesDeleteByIdCommand(a.exercisesService)
}
