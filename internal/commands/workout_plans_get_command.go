package commands

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
)

type WorkoutPlansGetCommand struct {
	AuthorizedCommand
	context             *CommandContext
	workoutPlansService *services.WorkoutPlansService
}

func NewWorkoutPlansGetCommandCommand(workoutPlanService *services.WorkoutPlansService) (*WorkoutPlansGetCommand, error) {
	context := &CommandContext{}
	return &WorkoutPlansGetCommand{context: context, workoutPlansService: workoutPlanService}, nil
}

func (c *WorkoutPlansGetCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *WorkoutPlansGetCommand) Validate() error {
	return nil
}

func (c *WorkoutPlansGetCommand) Execute() (interface{}, error) {
	workoutPlans, err := c.workoutPlansService.GetWorkoutPlans()
	if err != nil {
		return nil, err
	}
	return mapWorkoutPlansModelToShortDto(workoutPlans), nil
}

func mapWorkoutPlansModelToShortDto(workoutPlans []*models.WorkoutPlan) []*dto.WorkoutPlanShortModel {
	result := make([]*dto.WorkoutPlanShortModel, 0)
	for _, workoutPlan := range workoutPlans {
		result = append(result, &dto.WorkoutPlanShortModel{Id: workoutPlan.Id,
			Name:             workoutPlan.Name,
			ShortDescription: workoutPlan.ShortDescription,
			Repeatable:       workoutPlan.Repeatable})
	}
	return result
}
