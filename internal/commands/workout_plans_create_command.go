package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
)

type WorkoutPlansCreateCommand struct {
	AuthorizedCommand
	contextModel        *dto.WorkoutPlanCreateModel
	context             *CommandContext
	workoutPlansService *services.WorkoutPlansService
}

func NewWorkoutPlansCreateCommand(workoutPlanService *services.WorkoutPlansService) (*WorkoutPlansCreateCommand, error) {
	contextModel := &dto.WorkoutPlanCreateModel{}
	context := &CommandContext{CommandContent: contextModel}
	return &WorkoutPlansCreateCommand{context: context, workoutPlansService: workoutPlanService, contextModel: contextModel}, nil
}

func (c *WorkoutPlansCreateCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *WorkoutPlansCreateCommand) Validate() error {
	validationError := new(models.ValidationError)
	if c.contextModel.Name == "" {
		validationError.AddError("name", errors.New(models.ArgumentNullOrEmptyError))
	}
	if validationError.HasErrors() {
		return validationError
	}
	return nil
}

func (c *WorkoutPlansCreateCommand) Execute() (interface{}, error) {
	createModel := models.WorkoutPlan{Name: c.contextModel.Name,
		ShortDescription: c.contextModel.ShortDescription,
		Owner:            c.claims.UserId,
		Repeatable:       c.contextModel.Repeatable}
	for _, id := range c.contextModel.Workouts {
		item := new(models.Workout)
		item.Id = id
		createModel.Workouts = append(createModel.Workouts, item)
	}
	createdModel, err := c.workoutPlansService.CreateWorkoutPlan(createModel)
	if err != nil {
		return nil, err
	}
	return mapWorkoutPlanModelToDto(createdModel), nil
}

func mapWorkoutPlanModelToDto(workoutPlanModel *models.WorkoutPlan) *dto.WorkoutPlanFullModel {
	result := new(dto.WorkoutPlanFullModel)
	result.Id = workoutPlanModel.Id
	result.Name = workoutPlanModel.Name
	result.ShortDescription = workoutPlanModel.ShortDescription
	result.Repeatable = workoutPlanModel.Repeatable
	for _, item := range workoutPlanModel.Workouts {
		result.Workouts = append(result.Workouts, mapWorkoutModelToShortDto(item))
	}
	return result
}
