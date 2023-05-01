package services

import (
	"context"
	"github.com/go-logr/zapr"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

type UserWorkoutsCalendarGenerator struct {
	calendarService     *UserWorkoutsCalendarService
	userWorkoutsService *UserWorkoutsService
	logger              *zap.SugaredLogger
	cronForWorkouts     *cron.Cron
}

func NewUserWorkoutsCalendarGenerator(lc fx.Lifecycle,
	logger *zap.SugaredLogger,
	zapLoggerOriginal *zap.Logger,
	calendarService *UserWorkoutsCalendarService,
	userWorkoutsService *UserWorkoutsService) (*UserWorkoutsCalendarGenerator, error) {
	generator := &UserWorkoutsCalendarGenerator{logger: logger,
		calendarService:     calendarService,
		userWorkoutsService: userWorkoutsService}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := generator.generateWorkoutsOfDay()
			if err != nil {
				return err
			}

			cronLogger := zapr.NewLogger(zapLoggerOriginal)
			generator.cronForWorkouts = cron.New(cron.WithLogger(cronLogger))

			_, err = generator.cronForWorkouts.AddFunc("0 3 * * *", func() {
				generationError := generator.generateWorkoutsOfDay()
				if generationError != nil {
					return
				}
			})
			if err != nil {
				return err
			}

			logger.Info("Workout cron starting")
			generator.cronForWorkouts.Start()
			logger.Info("Workout cron started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Workout cron stopping")
			generator.cronForWorkouts.Stop()
			logger.Info("Workout cron stopped")
			return nil
		},
	})

	return generator, nil
}

func (g *UserWorkoutsCalendarGenerator) generateWorkoutsOfDay() error {
	g.logger.Info("Started generation of new calendar")
	// Get all active repeatable workouts
	activeWorkouts, err := g.userWorkoutsService.GetActiveRepeatableUserWorkouts()
	if err != nil {
		return err
	}
	var userWorkoutIds []string
	for _, item := range activeWorkouts {
		userWorkoutIds = append(userWorkoutIds, item.Id)
	}
	// Get last scheduled time for every active user workout
	lastSchedules, err := g.calendarService.
		workoutsStatisticRepository.GetLastWorkoutStatisticForUserWorkouts(userWorkoutIds)
	if err != nil {
		return err
	}
	_ = lastSchedules
	var newStatistic []*models.WorkoutStatistic
	// Add two week for generation
	twoWeek := utils.GetTodayUtc().AddDate(0, 0, 7*2)
	for _, item := range activeWorkouts {
		var lastSchedule time.Time
		var lastWorkout *models.Workout
		if value, ok := lastSchedules[item.Id]; ok && value != nil {
			lastSchedule = value.ScheduledDate
			lastSchedule = lastSchedule.AddDate(0, 0, 1)
			lastWorkout = getNextWorkout(item.WorkoutPlan.Workouts, value.Workout)
		} else {
			lastSchedule = utils.TruncateToDay(item.Created)
			lastWorkout = item.WorkoutPlan.Workouts[0]
		}

		// Generate a scheduled calendar for each next day
		for ; lastSchedule.Before(twoWeek); lastSchedule = lastSchedule.AddDate(0, 0, 1) {
			if weekdayContainsInSchedule(item.Schedule, lastSchedule.Weekday()) {
				statistic := &models.WorkoutStatistic{UserWorkout: &models.UserWorkout{}, Workout: &models.Workout{}}
				statistic.ScheduledDate = lastSchedule
				statistic.UserWorkout = item
				statistic.Workout = lastWorkout
				newStatistic = append(newStatistic, statistic)
				lastWorkout = getNextWorkout(item.WorkoutPlan.Workouts, lastWorkout)
			}
		}
	}
	if len(newStatistic) > 0 {
		err = g.calendarService.workoutsStatisticRepository.CreateWorkoutStatistics(newStatistic)
		if err != nil {
			return err
		}
	}
	g.logger.Infof("Finished generation of new calendar. New items created %d", len(newStatistic))
	return nil
}

func weekdayContainsInSchedule(schedule []time.Weekday, compareValue time.Weekday) bool {
	for _, item := range schedule {
		if item == compareValue {
			return true
		}
	}
	return false
}

func getNextWorkout(items []*models.Workout, current *models.Workout) *models.Workout {
	nextIndex := 0
	for i, item := range items {
		if item.Id == current.Id {
			nextIndex = i + 1
			break
		}
	}
	if nextIndex >= len(items) {
		nextIndex = 0
	}
	return items[nextIndex]
}
