package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
)

type ExercisesCreateCommand struct {
	AuthorizedCommand
	contextModel     *dto.ExerciseCreateModel
	context          *CommandContext
	exercisesService *services.ExercisesService
}

func NewExercisesCreateCommand(exerciseService *services.ExercisesService) (*ExercisesCreateCommand, error) {
	contextModel := &dto.ExerciseCreateModel{}
	context := &CommandContext{CommandContent: contextModel}
	return &ExercisesCreateCommand{context: context, exercisesService: exerciseService, contextModel: contextModel}, nil
}

func (c *ExercisesCreateCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *ExercisesCreateCommand) Validate() error {
	validationError := new(models.ValidationError)
	if c.contextModel.Name == "" {
		validationError.AddError("name", errors.New(models.ArgumentNullOrEmptyError))
	}
	if validationError.HasErrors() {
		return validationError
	}
	return nil
}

func (c *ExercisesCreateCommand) Execute() (interface{}, error) {
	createModel := models.Exercise{Name: c.contextModel.Name,
		ShortDescription: c.contextModel.ShortDescription,
		Owner:            c.claims.UserId}
	for _, complexId := range c.contextModel.Complex {
		exercise := new(models.Exercise)
		exercise.Id = complexId
		createModel.Complex = append(createModel.Complex, exercise)
	}
	createdModel, err := c.exercisesService.CreateExercise(createModel)
	if err != nil {
		return nil, err
	}
	return mapExerciseModelToDto(createdModel), nil
}

func mapExerciseModelToDto(exerciseModel *models.Exercise) *dto.ExerciseFullModel {
	result := new(dto.ExerciseFullModel)
	result.Id = exerciseModel.Id
	result.Name = exerciseModel.Name
	result.ShortDescription = exerciseModel.ShortDescription
	for _, complexItem := range exerciseModel.Complex {
		result.Complex = append(result.Complex, mapExerciseModelToDto(complexItem))
	}
	return result
}
