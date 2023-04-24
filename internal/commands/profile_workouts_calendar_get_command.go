package commands

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
)

type ProfileWorkoutsCalendarGetCommand struct {
	AuthorizedCommand
	context         *CommandContext
	calendarService *services.UserWorkoutsCalendarService
}

func NewProfileWorkoutsCalendarGetCommand(calendarService *services.UserWorkoutsCalendarService) (*ProfileWorkoutsCalendarGetCommand, error) {
	context := &CommandContext{}
	return &ProfileWorkoutsCalendarGetCommand{context: context, calendarService: calendarService}, nil
}

func (c *ProfileWorkoutsCalendarGetCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *ProfileWorkoutsCalendarGetCommand) Validate() error {
	return nil
}

func (c *ProfileWorkoutsCalendarGetCommand) Execute() (interface{}, error) {
	calendar, err := c.calendarService.GetCalendarForUser(c.claims.UserId)
	if err != nil {
		return nil, err
	}
	return mapWorkoutsCalendarToDto(calendar), nil
}

func mapWorkoutsCalendarToDto(calendar *models.WorkoutsCalendar) *dto.WorkoutsCalendarModel {
	result := new(dto.WorkoutsCalendarModel)
	for _, historyItem := range calendar.History {
		scheduledDate := dto.JsonDate{Time: historyItem.ScheduledDate}
		var workoutDate *dto.JsonDate
		if historyItem.WorkoutDate != nil {
			workoutDate = &dto.JsonDate{Time: *historyItem.WorkoutDate}
		}
		result.History = append(result.History, &dto.WorkoutStatisticModel{Id: historyItem.Id,
			ScheduledDate: scheduledDate, WorkoutDate: workoutDate})
	}

	if calendar.Current != nil {
		result.Current = new(dto.WorkoutStatisticModel)
		result.Current.Id = calendar.Current.Id
		result.Current.ScheduledDate = dto.JsonDate{Time: calendar.Current.ScheduledDate}
		if calendar.Current.WorkoutDate != nil {
			result.Current.WorkoutDate = &dto.JsonDate{Time: *calendar.Current.WorkoutDate}
		}
		result.Current.Workout = mapWorkoutModelToDto(calendar.Current.Workout)
	}

	for _, historyItem := range calendar.Upcoming {
		scheduledDate := dto.JsonDate{Time: historyItem.ScheduledDate}
		var workoutDate *dto.JsonDate
		if historyItem.WorkoutDate != nil {
			workoutDate = &dto.JsonDate{Time: *historyItem.WorkoutDate}
		}
		result.Upcoming = append(result.Upcoming, &dto.WorkoutStatisticModel{Id: historyItem.Id,
			ScheduledDate: scheduledDate, WorkoutDate: workoutDate})
	}

	return result
}
