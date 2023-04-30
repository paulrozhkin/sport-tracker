package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
	"time"
)

type ProfileWorkoutsCalendarVisitCommand struct {
	AuthorizedCommand
	workoutVisitId  string
	context         *CommandContext
	contextModel    *dto.ConfirmVisitModel
	calendarService *services.UserWorkoutsCalendarService
}

func NewProfileWorkoutsCalendarVisitCommand(calendarService *services.UserWorkoutsCalendarService) (*ProfileWorkoutsCalendarVisitCommand, error) {
	contextModel := new(dto.ConfirmVisitModel)
	context := &CommandContext{CommandContent: contextModel,
		CommandParameters: map[string]string{
			"workoutVisitId": "",
		}}
	return &ProfileWorkoutsCalendarVisitCommand{context: context, calendarService: calendarService, contextModel: contextModel}, nil
}

func (c *ProfileWorkoutsCalendarVisitCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *ProfileWorkoutsCalendarVisitCommand) Validate() error {
	c.workoutVisitId = c.context.CommandParameters["workoutVisitId"]
	validationError := new(models.ValidationError)
	if c.workoutVisitId == "" {
		validationError.AddError("workoutPlanId", errors.New(models.ArgumentNullOrEmptyError))
	}
	if !utils.IsValidUUID(c.workoutVisitId) {
		validationError.AddError("workoutPlanId", errors.New(models.InvalidFormat))
	}
	if validationError.HasErrors() {
		return validationError
	}
	return nil
}

func (c *ProfileWorkoutsCalendarVisitCommand) Execute() (interface{}, error) {
	var dateWithoutTime *time.Time
	if c.contextModel.WorkoutDate != nil {
		date := utils.TruncateToDay(c.contextModel.WorkoutDate.Time)
		dateWithoutTime = &date
	}
	visitModel := models.ConfirmVisit{WorkoutVisitId: c.workoutVisitId, Comment: c.contextModel.Comment,
		WorkoutDate: dateWithoutTime}
	result, err := c.calendarService.ConfirmVisit(visitModel)
	if err != nil {
		return nil, err
	}
	return mapWorkoutStatisticToDto(result), nil
}
