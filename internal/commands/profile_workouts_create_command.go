package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type ProfileWorkoutsCreateCommand struct {
	AuthorizedCommand
	contextModel        *dto.ProfileWorkoutCreateModel
	context             *CommandContext
	userWorkoutsService *services.UserWorkoutsService
}

func NewProfileWorkoutsCreateCommand(userWorkoutsService *services.UserWorkoutsService) (*ProfileWorkoutsCreateCommand, error) {
	contextModel := &dto.ProfileWorkoutCreateModel{}
	context := &CommandContext{CommandContent: contextModel}
	return &ProfileWorkoutsCreateCommand{context: context, userWorkoutsService: userWorkoutsService, contextModel: contextModel}, nil
}

func (c *ProfileWorkoutsCreateCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *ProfileWorkoutsCreateCommand) Validate() error {
	validationError := new(models.ValidationError)
	if c.contextModel.WorkoutPlan == "" {
		validationError.AddError("workoutPlan", errors.New(models.ArgumentNullOrEmptyError))
	}
	if !utils.IsValidUUID(c.contextModel.WorkoutPlan) {
		validationError.AddError("workoutPlan", errors.New(models.InvalidFormat))
	}
	if c.contextModel.Schedule == nil {
		validationError.AddError("workoutPlan", errors.New(models.ArgumentNullOrEmptyError))
	} else if len(c.contextModel.Schedule) == 0 {
		validationError.AddError("workoutPlan", errors.New(models.ArgumentNullOrEmptyError))
	}
	if validationError.HasErrors() {
		return validationError
	}
	return nil
}

func (c *ProfileWorkoutsCreateCommand) Execute() (interface{}, error) {
	createModel := models.UserWorkout{UserId: c.claims.UserId,
		WorkoutPlan: &models.WorkoutPlan{},
		Schedule:    mapScheduleDtoToModel(c.contextModel.Schedule)}
	createModel.WorkoutPlan.Id = c.contextModel.WorkoutPlan
	createdModel, err := c.userWorkoutsService.CreateUserWorkout(createModel)
	if err != nil {
		return nil, err
	}
	return mapProfileWorkoutModelToDto(createdModel), nil
}

func mapScheduleDtoToModel(dto []dto.DaysOfWeekDto) []models.DaysOfWeek {
	var result []models.DaysOfWeek
	mapUniq := make(map[models.DaysOfWeek]struct{})
	for _, item := range dto {
		modelValue := models.DaysOfWeek(item)
		mapUniq[modelValue] = struct{}{}
	}
	for key := range mapUniq {
		result = append(result, key)
	}
	return result
}

func mapProfileWorkoutModelToDto(profileWorkoutModel *models.UserWorkout) *dto.ProfileWorkoutShortModel {
	result := new(dto.ProfileWorkoutShortModel)
	result.Id = profileWorkoutModel.Id
	return result
}
