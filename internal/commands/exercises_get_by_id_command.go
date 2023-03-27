package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type ExercisesGetByIdCommand struct {
	AuthorizedCommand
	exerciseId       string
	context          *CommandContext
	exercisesService *services.ExercisesService
}

func NewExercisesGetByIdCommand(exerciseService *services.ExercisesService) (*ExercisesGetByIdCommand, error) {
	command := new(ExercisesGetByIdCommand)
	command.context = &CommandContext{CommandParameters: map[string]string{
		"exerciseId": "",
	}}
	command.exercisesService = exerciseService
	return command, nil
}

func (c *ExercisesGetByIdCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *ExercisesGetByIdCommand) Validate() error {
	c.exerciseId = c.context.CommandParameters["exerciseId"]
	validationError := new(models.ValidationError)
	if c.exerciseId == "" {
		validationError.AddError("exerciseId", errors.New(models.ArgumentNullOrEmptyError))
	}
	if !utils.IsValidUUID(c.exerciseId) {
		validationError.AddError("exerciseId", errors.New(models.InvalidFormat))
	}
	if validationError.HasErrors() {
		return validationError
	}
	return nil
}

func (c *ExercisesGetByIdCommand) Execute() (interface{}, error) {
	exercise, err := c.exercisesService.GetExerciseById(c.exerciseId)
	if err != nil {
		return nil, err
	}
	return mapExerciseModelToDto(exercise), nil
}
