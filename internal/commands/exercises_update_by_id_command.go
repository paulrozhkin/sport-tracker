package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type ExercisesUpdateByIdCommand struct {
	AuthorizedCommand
	contextModel     *dto.ExerciseCreateModel
	context          *CommandContext
	exercisesService *services.ExercisesService
	exerciseId       string
}

func NewExercisesUpdateByIdCommand(exerciseService *services.ExercisesService) (*ExercisesUpdateByIdCommand, error) {
	command := new(ExercisesUpdateByIdCommand)
	command.contextModel = &dto.ExerciseCreateModel{}
	command.context = &CommandContext{CommandParameters: map[string]string{
		"exerciseId": "",
	}, CommandContent: command.contextModel}
	command.exercisesService = exerciseService
	return command, nil
}

func (c *ExercisesUpdateByIdCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *ExercisesUpdateByIdCommand) Validate() error {
	c.exerciseId = c.context.CommandParameters["exerciseId"]
	validationError := new(models.ValidationError)
	if c.contextModel.Name == "" {
		validationError.AddError("name", errors.New(models.ArgumentNullOrEmptyError))
	}
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

func (c *ExercisesUpdateByIdCommand) Execute() (interface{}, error) {
	model := models.Exercise{
		Name:             c.contextModel.Name,
		ShortDescription: c.contextModel.ShortDescription,
		Owner:            c.claims.UserId,
	}
	model.Id = c.exerciseId
	for _, complexId := range c.contextModel.Complex {
		exercise := new(models.Exercise)
		exercise.Id = complexId
		model.Complex = append(model.Complex, exercise)
	}
	updatedModel, err := c.exercisesService.UpdateExercise(model)
	if err != nil {
		return nil, err
	}
	return mapExerciseModelToDto(updatedModel), nil
}
