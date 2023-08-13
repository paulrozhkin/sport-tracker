package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"go.uber.org/zap"
	"net/http"
)

type HealthcheckRoute struct {
	store  *infrastructure.Store
	logger *zap.SugaredLogger
}

func NewHealthcheckRoute(store *infrastructure.Store,
	logger *zap.SugaredLogger) *HealthcheckRoute {
	return &HealthcheckRoute{store: store, logger: logger}
}

func (*HealthcheckRoute) Method() string {
	return http.MethodGet
}

func (*HealthcheckRoute) Pattern() string {
	return "/healthcheck"
}

func (r *HealthcheckRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewHealthcheckCommand(r.store, r.logger)
}
