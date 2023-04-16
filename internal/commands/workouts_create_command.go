package commands

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
)

type WorkoutsCreateCommand struct {
	AuthorizedCommand
	contextModel    *dto.WorkoutCreateModel
	context         *CommandContext
	workoutsService *services.WorkoutsService
}

func NewWorkoutsCreateCommand(workoutsService *services.WorkoutsService) (*WorkoutsCreateCommand, error) {
	contextModel := &dto.WorkoutCreateModel{}
	context := &CommandContext{CommandContent: contextModel}
	return &WorkoutsCreateCommand{context: context, workoutsService: workoutsService, contextModel: contextModel}, nil
}

func (c *WorkoutsCreateCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *WorkoutsCreateCommand) Validate() error {
	validationError := new(models.ValidationError)
	if validationError.HasErrors() {
		return validationError
	}
	return nil
}

func (c *WorkoutsCreateCommand) Execute() (interface{}, error) {
	createModel := models.Workout{CustomName: c.contextModel.CustomName,
		CustomDescription: c.contextModel.CustomDescription,
		Owner:             c.claims.UserId}
	for _, complexId := range c.contextModel.Complex {
		exercise := new(models.Exercise)
		exercise.Id = complexId
		createModel.Complex = append(createModel.Complex, exercise)
	}
	createdModel, err := c.workoutsService.CreateWorkout(createModel)
	if err != nil {
		return nil, err
	}
	return mapWorkoutModelToDto(createdModel), nil
}

func mapWorkoutModelToDto(exerciseModel *models.Workout) *dto.WorkoutFullModel {
	result := new(dto.WorkoutFullModel)
	result.Id = exerciseModel.Id
	result.CustomName = exerciseModel.CustomName
	result.CustomDescription = exerciseModel.CustomDescription
	for _, complexItem := range exerciseModel.Complex {
		result.Complex = append(result.Complex, mapExerciseModelToDto(complexItem))
	}
	return result
}

func mapWorkoutModelToShortDto(exerciseModel *models.Workout) *dto.WorkoutShortModel {
	result := new(dto.WorkoutShortModel)
	result.Id = exerciseModel.Id
	result.CustomName = exerciseModel.CustomName
	result.CustomDescription = exerciseModel.CustomDescription
	return result
}
