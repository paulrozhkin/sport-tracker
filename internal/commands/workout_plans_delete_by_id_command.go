package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type WorkoutPlansDeleteByIdCommand struct {
	AuthorizedCommand
	workoutPlanId       string
	context             *CommandContext
	workoutPlansService *services.WorkoutPlansService
}

func NewWorkoutPlansDeleteByIdCommand(workoutPlanService *services.WorkoutPlansService) (*WorkoutPlansDeleteByIdCommand, error) {
	command := new(WorkoutPlansDeleteByIdCommand)
	command.context = &CommandContext{CommandParameters: map[string]string{
		"workoutPlanId": "",
	}}
	command.workoutPlansService = workoutPlanService
	return command, nil
}

func (c *WorkoutPlansDeleteByIdCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *WorkoutPlansDeleteByIdCommand) Validate() error {
	c.workoutPlanId = c.context.CommandParameters["workoutPlanId"]
	validationError := new(models.ValidationError)
	if c.workoutPlanId == "" {
		validationError.AddError("workoutPlanId", errors.New(models.ArgumentNullOrEmptyError))
	}
	if !utils.IsValidUUID(c.workoutPlanId) {
		validationError.AddError("workoutPlanId", errors.New(models.InvalidFormat))
	}
	if validationError.HasErrors() {
		return validationError
	}
	return nil
}

func (c *WorkoutPlansDeleteByIdCommand) Execute() (interface{}, error) {
	return nil, c.workoutPlansService.DeleteWorkoutPlanById(c.workoutPlanId)
}
