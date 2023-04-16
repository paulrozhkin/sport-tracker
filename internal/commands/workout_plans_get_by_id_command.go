package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type WorkoutPlansGetByIdCommand struct {
	AuthorizedCommand
	workoutPlanId       string
	context             *CommandContext
	workoutPlansService *services.WorkoutPlansService
}

func NewWorkoutPlansGetByIdCommand(workoutPlanService *services.WorkoutPlansService) (*WorkoutPlansGetByIdCommand, error) {
	command := new(WorkoutPlansGetByIdCommand)
	command.context = &CommandContext{CommandParameters: map[string]string{
		"workoutPlanId": "",
	}}
	command.workoutPlansService = workoutPlanService
	return command, nil
}

func (c *WorkoutPlansGetByIdCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *WorkoutPlansGetByIdCommand) Validate() error {
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

func (c *WorkoutPlansGetByIdCommand) Execute() (interface{}, error) {
	workoutPlan, err := c.workoutPlansService.GetWorkoutPlanById(c.workoutPlanId)
	if err != nil {
		return nil, err
	}
	return mapWorkoutPlanModelToDto(workoutPlan), nil
}
