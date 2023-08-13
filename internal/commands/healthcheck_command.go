package commands

import (
	"context"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"go.uber.org/zap"
	"time"
)

type HealthcheckCommand struct {
	UnauthorizedCommand
	context *CommandContext
	store   *infrastructure.Store
	logger  *zap.SugaredLogger
}

func NewHealthcheckCommand(store *infrastructure.Store,
	logger *zap.SugaredLogger) (*HealthcheckCommand, error) {
	context := &CommandContext{}
	return &HealthcheckCommand{context: context, store: store, logger: logger}, nil
}

func (c *HealthcheckCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *HealthcheckCommand) Validate() error {
	return nil
}

func (c *HealthcheckCommand) Execute() (interface{}, error) {
	timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	conn, err := c.store.Pool.Acquire(timeout)
	if err != nil {
		c.logger.Errorf("healthcheck failed while get connection due to: %v", err)
		return nil, err
	}
	errPing := conn.Ping(timeout)
	conn.Release()
	if errPing != nil {
		c.logger.Errorf("healthcheck failed while ping connection due to: %v", errPing)
		return nil, errPing
	}
	stat := c.store.Pool.Stat()

	result := &dto.Healthcheck{
		TotalDbInvokes:      stat.AcquireCount(),
		CurrentDbConnection: int(stat.AcquiredConns()),
		MaxDbConnections:    int(stat.MaxConns()),
	}
	return result, nil
}
