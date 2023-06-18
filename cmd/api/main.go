package main

import (
	"github.com/paulrozhkin/sport-tracker/config"
	"github.com/paulrozhkin/sport-tracker/internal/http_server"
	"github.com/paulrozhkin/sport-tracker/internal/http_server/routes"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"github.com/paulrozhkin/sport-tracker/internal/metrics"
	"github.com/paulrozhkin/sport-tracker/internal/repositories"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	fx.New(
		fx.Provide(
			config.InitLogger,
			zap.L,
			zap.S,
			config.LoadConfigurations,
			http_server.NewHTTPServer,
			infrastructure.NewPostgresStore,
			services.NewDataSeedingService,
			fx.Annotate(
				http_server.NewServerRoute,
				fx.ParamTags(`group:"routes"`),
			),
		),
		createServicesRegistration(),
		createRoutesRegistration(),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Invoke(func(*config.LoggerConfigurator) {}),
		fx.Invoke(func(*infrastructure.Store) {}),
		fx.Invoke(func(seeding *services.DataSeedingService) {}),
		fx.Invoke(func(*services.UserWorkoutsCalendarGenerator) {}),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

func createRoutesRegistration() fx.Option {
	return fx.Provide(
		routes.AsRoute(routes.NewAuthRoute),
		routes.AsRoute(routes.NewRegisterRoute),
		// Exercises
		routes.AsRoute(routes.NewExercisesCreateRoute),
		routes.AsRoute(routes.NewExercisesGetRoute),
		routes.AsRoute(routes.NewExercisesGetByIdRoute),
		routes.AsRoute(routes.NewExercisesDeleteByIdRoute),
		routes.AsRoute(routes.NewExercisesUpdateByIdRoute),
		// Workouts
		routes.AsRoute(routes.NewWorkoutsCreateRoute),
		routes.AsRoute(routes.NewWorkoutsGetRoute),
		routes.AsRoute(routes.NewWorkoutsGetByIdRoute),
		routes.AsRoute(routes.NewWorkoutsUpdateByIdRoute),
		routes.AsRoute(routes.NewWorkoutsDeleteByIdRoute),
		// Workout plans
		routes.AsRoute(routes.NewWorkoutPlansCreateRoute),
		routes.AsRoute(routes.NewWorkoutPlansGetRoute),
		routes.AsRoute(routes.NewWorkoutPlansGetByIdRoute),
		routes.AsRoute(routes.NewWorkoutPlansUpdateByIdRoute),
		routes.AsRoute(routes.NewWorkoutPlansDeleteByIdRoute),
		// Profile workouts
		routes.AsRoute(routes.NewProfileWorkoutsCreateRoute),
		routes.AsRoute(routes.NewProfileWorkoutsDeleteRoute),
		// Profile workouts calendar
		routes.AsRoute(routes.NewProfileWorkoutsCalendarGetRoute),
		routes.AsRoute(routes.NewProfileWorkoutsCalendarVisitRoute),
		routes.AsRoute(routes.NewProfileWorkoutsCalendarGetByIdRoute),
		// Statistic
		routes.AsRoute(routes.NewProfileStatisticGetRoute),
		routes.AsRoute(routes.NewProfileStatisticWeightCreateRoute),
		routes.AsRoute(routes.NewProfileStatisticWeightUpdateByIdRoute),
		routes.AsRoute(routes.NewProfileStatisticWeightDeleteByIdRoute),
	)
}

func createServicesRegistration() fx.Option {
	return fx.Provide(
		services.NewUsersService,
		repositories.NewUsersRepository,
		services.NewTokenService,
		repositories.NewExercisesRepository,
		services.NewExercisesService,
		services.NewWorkoutsService,
		repositories.NewWorkoutsRepository,
		services.NewWorkoutPlansService,
		repositories.NewWorkoutPlansRepository,
		services.NewUserWorkoutsService,
		repositories.NewUserWorkoutsRepository,
		services.NewUserWorkoutsCalendarService,
		repositories.NewWorkoutsStatisticRepository,
		services.NewUserWorkoutsCalendarGenerator,
		repositories.NewUserWeightMeasurementRepository,
		services.NewUserStatisticService,
		metrics.NewTrafficMetrics,
		http_server.NewTrafficMiddleware,
		metrics.NewUsersMetrics,
	)
}
