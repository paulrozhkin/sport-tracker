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
	"log"
	"net/http"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
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
	)
}

func createServicesRegistration() fx.Option {
	return fx.Provide(
		services.NewUserService,
		repositories.NewUsersRepository,
	)
}
