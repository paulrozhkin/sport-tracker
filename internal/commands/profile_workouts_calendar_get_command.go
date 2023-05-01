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
	result.History = make([]*dto.WorkoutStatisticModel, 0)
	result.Upcoming = make([]*dto.WorkoutStatisticModel, 0)
	for _, historyItem := range calendar.History {
		result.History = append(result.History, mapWorkoutStatisticToShortDto(historyItem))
	}

	if calendar.Current != nil {
		result.Current = mapWorkoutStatisticToDto(calendar.Current)
	}

	for _, historyItem := range calendar.Upcoming {
		result.Upcoming = append(result.Upcoming, mapWorkoutStatisticToShortDto(historyItem))
	}

	return result
}

func mapWorkoutStatisticToDto(workoutStatistic *models.WorkoutStatistic) *dto.WorkoutStatisticModel {
	result := new(dto.WorkoutStatisticModel)
	result.Id = workoutStatistic.Id
	result.ScheduledDate = dto.JsonDate{Time: workoutStatistic.ScheduledDate}
	if workoutStatistic.WorkoutDate != nil {
		result.WorkoutDate = &dto.JsonDate{Time: *workoutStatistic.WorkoutDate}
	}
	result.Comment = workoutStatistic.Comment
	if workoutStatistic.Workout != nil {
		result.Workout = mapWorkoutModelToDto(workoutStatistic.Workout)
	}
	return result
}

func mapWorkoutStatisticToShortDto(workoutStatistic *models.WorkoutStatistic) *dto.WorkoutStatisticModel {
	result := new(dto.WorkoutStatisticModel)
	result.Id = workoutStatistic.Id
	result.ScheduledDate = dto.JsonDate{Time: workoutStatistic.ScheduledDate}
	if workoutStatistic.WorkoutDate != nil {
		result.WorkoutDate = &dto.JsonDate{Time: *workoutStatistic.WorkoutDate}
	}
	return result
}
