package main

import (
	"github.com/paulrozhkin/sport-tracker/config"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
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
			NewHTTPServer,
			NewServerRoute,
			infrastructure.NewStore,
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Invoke(func(*infrastructure.Store) {}),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
