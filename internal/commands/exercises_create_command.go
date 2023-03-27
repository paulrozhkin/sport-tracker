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

func (a *ExercisesCreateCommand) GetCommandContext() *CommandContext {
	return a.context
}

func (a *ExercisesCreateCommand) Validate() error {
	validationError := new(models.ValidationError)
	if a.contextModel.Name == "" {
		validationError.AddError("name", errors.New(models.ArgumentNullOrEmptyError))
	}
	if validationError.HasErrors() {
		return validationError
	}
	return nil
}

func (a *ExercisesCreateCommand) Execute() (interface{}, error) {
	createModel := models.Exercise{Name: a.contextModel.Name,
		ShortDescription: a.contextModel.ShortDescription,
		Owner:            a.claims.UserId}
	for _, complexId := range a.contextModel.Complex {
		exercise := new(models.Exercise)
		exercise.Id = complexId
		createModel.Complex = append(createModel.Complex, exercise)
	}
	createdModel, err := a.exercisesService.CreateExercise(createModel)
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
