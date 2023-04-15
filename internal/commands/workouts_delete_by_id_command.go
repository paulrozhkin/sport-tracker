package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type WorkoutsDeleteByIdCommand struct {
	AuthorizedCommand
	workoutId       string
	context         *CommandContext
	workoutsService *services.WorkoutsService
}

func NewWorkoutsDeleteByIdCommand(workoutService *services.WorkoutsService) (*WorkoutsDeleteByIdCommand, error) {
	command := new(WorkoutsDeleteByIdCommand)
	command.context = &CommandContext{CommandParameters: map[string]string{
		"workoutId": "",
	}}
	command.workoutsService = workoutService
	return command, nil
}

func (c *WorkoutsDeleteByIdCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *WorkoutsDeleteByIdCommand) Validate() error {
	c.workoutId = c.context.CommandParameters["workoutId"]
	validationError := new(models.ValidationError)
	if c.workoutId == "" {
		validationError.AddError("workoutId", errors.New(models.ArgumentNullOrEmptyError))
	}
	if !utils.IsValidUUID(c.workoutId) {
		validationError.AddError("workoutId", errors.New(models.InvalidFormat))
	}
	if validationError.HasErrors() {
		return validationError
	}
	return nil
}

func (c *WorkoutsDeleteByIdCommand) Execute() (interface{}, error) {
	return nil, c.workoutsService.DeleteWorkoutById(c.workoutId)
}
