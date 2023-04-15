package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type WorkoutsGetByIdCommand struct {
	AuthorizedCommand
	workoutId       string
	context         *CommandContext
	workoutsService *services.WorkoutsService
}

func NewWorkoutsGetByIdCommand(workoutsService *services.WorkoutsService) (*WorkoutsGetByIdCommand, error) {
	command := new(WorkoutsGetByIdCommand)
	command.context = &CommandContext{CommandParameters: map[string]string{
		"workoutId": "",
	}}
	command.workoutsService = workoutsService
	return command, nil
}

func (c *WorkoutsGetByIdCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *WorkoutsGetByIdCommand) Validate() error {
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

func (c *WorkoutsGetByIdCommand) Execute() (interface{}, error) {
	workout, err := c.workoutsService.GetWorkoutById(c.workoutId)
	if err != nil {
		return nil, err
	}
	return mapWorkoutModelToDto(workout), nil
}
