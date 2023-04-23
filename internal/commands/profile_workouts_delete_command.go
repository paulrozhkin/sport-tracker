package commands

import (
	"github.com/paulrozhkin/sport-tracker/internal/services"
)

type ProfileWorkoutsDeleteCommand struct {
	AuthorizedCommand
	context             *CommandContext
	userWorkoutsService *services.UserWorkoutsService
}

func NewProfileWorkoutsDeleteCommand(userWorkoutsService *services.UserWorkoutsService) (*ProfileWorkoutsDeleteCommand, error) {
	context := &CommandContext{}
	return &ProfileWorkoutsDeleteCommand{context: context, userWorkoutsService: userWorkoutsService}, nil
}

func (c *ProfileWorkoutsDeleteCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *ProfileWorkoutsDeleteCommand) Validate() error {
	return nil
}

func (c *ProfileWorkoutsDeleteCommand) Execute() (interface{}, error) {
	_, err := c.userWorkoutsService.DeactivateWorkoutForUser(c.claims.UserId)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
