package main

import (
	"github.com/paulrozhkin/sport-tracker/config"
	"github.com/paulrozhkin/sport-tracker/internal/http_server"
	"github.com/paulrozhkin/sport-tracker/internal/http_server/routes"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"github.com/paulrozhkin/sport-tracker/internal/repositories"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	// TODO: refactoring with config from file
	production := false
	var logger *zap.Logger
	if production {
		loggerCfg := zap.NewProductionConfig()
		//loggerCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		loggerCfg.OutputPaths = append(loggerCfg.OutputPaths, "./logs/sport-tracker.log")
		logger = zap.Must(loggerCfg.Build())
	} else {
		logger, _ = zap.NewDevelopment()
	}

	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	fx.New(
		fx.Provide(
			zap.L,
			zap.S,
			config.LoadConfigurations,
			http_server.NewHTTPServer,
			infrastructure.NewPostgresStore,
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
		fx.Invoke(func(*infrastructure.Store) {}),
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
	)
}
