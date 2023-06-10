package config

import (
	"context"
	"encoding/json"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type LoggerConfigurator struct {
	logger     *zap.Logger
	undoLogger func()
}

func InitLogger(lc fx.Lifecycle,
	cfg *Configuration) (*LoggerConfigurator, error) {
	configurator := new(LoggerConfigurator)
	if cfg.Production {
		configurator.logger, _ = zap.NewProduction()
	} else {
		configurator.logger, _ = zap.NewDevelopment()
	}
	configurator.undoLogger = zap.ReplaceGlobals(configurator.logger)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			confString, _ := json.MarshalIndent(cfg, "", " ")
			configurator.logger.Sugar().Info("Configuration:\n", string(confString))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_ = configurator.logger.Sync()
			configurator.undoLogger()
			return nil
		},
	})
	return configurator, nil
}
