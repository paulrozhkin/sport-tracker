package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type ExercisesDeleteByIdCommand struct {
	AuthorizedCommand
	exerciseId       string
	context          *CommandContext
	exercisesService *services.ExercisesService
}

func NewExercisesDeleteByIdCommand(exerciseService *services.ExercisesService) (*ExercisesDeleteByIdCommand, error) {
	command := new(ExercisesDeleteByIdCommand)
	command.context = &CommandContext{CommandParameters: map[string]string{
		"exerciseId": "",
	}}
	command.exercisesService = exerciseService
	return command, nil
}

func (c *ExercisesDeleteByIdCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *ExercisesDeleteByIdCommand) Validate() error {
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

func (c *ExercisesDeleteByIdCommand) Execute() (interface{}, error) {
	return nil, c.exercisesService.DeleteExerciseById(c.exerciseId)
}
