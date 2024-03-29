package commands

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
)

type ExercisesGetCommand struct {
	AuthorizedCommand
	context          *CommandContext
	exercisesService *services.ExercisesService
}

func NewExercisesGetCommandCommand(exerciseService *services.ExercisesService) (*ExercisesGetCommand, error) {
	context := &CommandContext{}
	return &ExercisesGetCommand{context: context, exercisesService: exerciseService}, nil
}

func (c *ExercisesGetCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *ExercisesGetCommand) Validate() error {
	return nil
}

func (c *ExercisesGetCommand) Execute() (interface{}, error) {
	exercises, err := c.exercisesService.GetExercises()
	if err != nil {
		return nil, err
	}
	return mapExercisesModelToShortDto(exercises), nil
}

func mapExercisesModelToShortDto(exercises []*models.Exercise) []*dto.ExerciseShortModel {
	result := make([]*dto.ExerciseShortModel, 0)
	for _, exercise := range exercises {
		result = append(result, &dto.ExerciseShortModel{Id: exercise.Id,
			Name:             exercise.Name,
			ShortDescription: exercise.ShortDescription})
	}
	return result
}
