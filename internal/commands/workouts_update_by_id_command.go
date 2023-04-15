package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type WorkoutsUpdateByIdCommand struct {
	AuthorizedCommand
	contextModel    *dto.WorkoutCreateModel
	context         *CommandContext
	workoutsService *services.WorkoutsService
	workoutId       string
}

func NewWorkoutsUpdateByIdCommand(workoutService *services.WorkoutsService) (*WorkoutsUpdateByIdCommand, error) {
	command := new(WorkoutsUpdateByIdCommand)
	command.contextModel = &dto.WorkoutCreateModel{}
	command.context = &CommandContext{CommandParameters: map[string]string{
		"workoutId": "",
	}, CommandContent: command.contextModel}
	command.workoutsService = workoutService
	return command, nil
}

func (c *WorkoutsUpdateByIdCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *WorkoutsUpdateByIdCommand) Validate() error {
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

func (c *WorkoutsUpdateByIdCommand) Execute() (interface{}, error) {
	model := models.Workout{
		CustomName:        c.contextModel.CustomName,
		CustomDescription: c.contextModel.CustomDescription,
		Owner:             c.claims.UserId,
	}
	model.Id = c.workoutId
	for _, complexId := range c.contextModel.Complex {
		exercise := new(models.Exercise)
		exercise.Id = complexId
		model.Complex = append(model.Complex, exercise)
	}
	updatedModel, err := c.workoutsService.UpdateWorkout(model)
	if err != nil {
		return nil, err
	}
	return mapWorkoutModelToDto(updatedModel), nil
}
