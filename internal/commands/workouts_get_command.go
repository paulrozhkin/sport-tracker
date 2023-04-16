package commands

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
)

type WorkoutsGetCommand struct {
	AuthorizedCommand
	context         *CommandContext
	workoutsService *services.WorkoutsService
}

func NewWorkoutsGetCommandCommand(workoutService *services.WorkoutsService) (*WorkoutsGetCommand, error) {
	context := &CommandContext{}
	return &WorkoutsGetCommand{context: context, workoutsService: workoutService}, nil
}

func (c *WorkoutsGetCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *WorkoutsGetCommand) Validate() error {
	return nil
}

func (c *WorkoutsGetCommand) Execute() (interface{}, error) {
	workouts, err := c.workoutsService.GetWorkouts()
	if err != nil {
		return nil, err
	}
	return mapWorkoutsModelToShortDto(workouts), nil
}

func mapWorkoutsModelToShortDto(workouts []*models.Workout) []*dto.WorkoutShortModel {
	result := make([]*dto.WorkoutShortModel, 0)
	for _, workout := range workouts {
		result = append(result, &dto.WorkoutShortModel{Id: workout.Id,
			CustomName:        workout.CustomName,
			CustomDescription: workout.CustomDescription})
	}
	return result
}
