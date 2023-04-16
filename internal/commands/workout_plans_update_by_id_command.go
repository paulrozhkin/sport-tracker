package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type WorkoutPlansUpdateByIdCommand struct {
	AuthorizedCommand
	contextModel        *dto.WorkoutPlanCreateModel
	context             *CommandContext
	workoutPlansService *services.WorkoutPlansService
	workoutPlanId       string
}

func NewWorkoutPlansUpdateByIdCommand(workoutPlanService *services.WorkoutPlansService) (*WorkoutPlansUpdateByIdCommand, error) {
	command := new(WorkoutPlansUpdateByIdCommand)
	command.contextModel = &dto.WorkoutPlanCreateModel{}
	command.context = &CommandContext{CommandParameters: map[string]string{
		"workoutPlanId": "",
	}, CommandContent: command.contextModel}
	command.workoutPlansService = workoutPlanService
	return command, nil
}

func (c *WorkoutPlansUpdateByIdCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *WorkoutPlansUpdateByIdCommand) Validate() error {
	c.workoutPlanId = c.context.CommandParameters["workoutPlanId"]
	validationError := new(models.ValidationError)
	if c.contextModel.Name == "" {
		validationError.AddError("name", errors.New(models.ArgumentNullOrEmptyError))
	}
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

func (c *WorkoutPlansUpdateByIdCommand) Execute() (interface{}, error) {
	model := models.WorkoutPlan{
		Name:             c.contextModel.Name,
		ShortDescription: c.contextModel.ShortDescription,
		Owner:            c.claims.UserId,
		Repeatable:       c.contextModel.Repeatable,
	}
	model.Id = c.workoutPlanId
	for _, id := range c.contextModel.Workouts {
		item := new(models.Workout)
		item.Id = id
		model.Workouts = append(model.Workouts, item)
	}
	updatedModel, err := c.workoutPlansService.UpdateWorkoutPlan(model)
	if err != nil {
		return nil, err
	}
	return mapWorkoutPlanModelToDto(updatedModel), nil
}
