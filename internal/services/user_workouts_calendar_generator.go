package services

import (
	"context"
	"github.com/go-logr/zapr"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type UserWorkoutsCalendarGenerator struct {
	calendarService *UserWorkoutsCalendarService
	logger          *zap.SugaredLogger
	cronForWorkouts *cron.Cron
}

func NewUserWorkoutsCalendarGenerator(lc fx.Lifecycle,
	logger *zap.SugaredLogger,
	zapLoggerOriginal *zap.Logger,
	calendarService *UserWorkoutsCalendarService) (*UserWorkoutsCalendarGenerator, error) {
	generator := &UserWorkoutsCalendarGenerator{logger: logger,
		calendarService: calendarService}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := calendarService.CreateCalendarForRepeatableActiveUserWorkouts()
			if err != nil {
				return err
			}

			cronLogger := zapr.NewLogger(zapLoggerOriginal)
			generator.cronForWorkouts = cron.New(cron.WithLogger(cronLogger))

			_, err = generator.cronForWorkouts.AddFunc("0 3 * * *", func() {
				generationError := calendarService.CreateCalendarForRepeatableActiveUserWorkouts()
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
