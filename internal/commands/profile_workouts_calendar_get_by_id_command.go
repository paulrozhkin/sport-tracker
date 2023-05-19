package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type ProfileWorkoutsCalendarGetByIdCommand struct {
	AuthorizedCommand
	workoutStatisticId   string
	context              *CommandContext
	userWorkoutsCalendar *services.UserWorkoutsCalendarService
}

func NewProfileWorkoutsCalendarGetByIdCommand(userWorkoutsCalendar *services.UserWorkoutsCalendarService) (*ProfileWorkoutsCalendarGetByIdCommand, error) {
	command := new(ProfileWorkoutsCalendarGetByIdCommand)
	command.context = &CommandContext{CommandParameters: map[string]string{
		"workoutStatisticId": "",
	}}
	command.userWorkoutsCalendar = userWorkoutsCalendar
	return command, nil
}

func (c *ProfileWorkoutsCalendarGetByIdCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *ProfileWorkoutsCalendarGetByIdCommand) Validate() error {
	c.workoutStatisticId = c.context.CommandParameters["workoutStatisticId"]
	validationError := new(models.ValidationError)
	if c.workoutStatisticId == "" {
		validationError.AddError("workoutStatisticId", errors.New(models.ArgumentNullOrEmptyError))
	}
	if !utils.IsValidUUID(c.workoutStatisticId) {
		validationError.AddError("workoutStatisticId", errors.New(models.InvalidFormat))
	}
	if validationError.HasErrors() {
		return validationError
	}
	return nil
}

func (c *ProfileWorkoutsCalendarGetByIdCommand) Execute() (interface{}, error) {
	workoutStatistic, err := c.userWorkoutsCalendar.GetWorkoutStatisticById(c.workoutStatisticId)
	if err != nil {
		return nil, err
	}
	return mapWorkoutStatisticToDto(workoutStatistic), nil
}
