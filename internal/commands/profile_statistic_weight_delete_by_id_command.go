package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type ProfileStatisticWeightDeleteByIdCommand struct {
	AuthorizedCommand
	context                  *CommandContext
	userStatisticService     *services.UserStatisticService
	profileStatisticWeightId string
}

func NewProfileStatisticWeightDeleteByIdCommand(userStatisticService *services.UserStatisticService) (*ProfileStatisticWeightDeleteByIdCommand, error) {
	context := &CommandContext{CommandParameters: map[string]string{
		"profileStatisticWeightId": "",
	}}
	return &ProfileStatisticWeightDeleteByIdCommand{context: context, userStatisticService: userStatisticService}, nil
}

func (c *ProfileStatisticWeightDeleteByIdCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *ProfileStatisticWeightDeleteByIdCommand) Validate() error {
	c.profileStatisticWeightId = c.context.CommandParameters["profileStatisticWeightId"]
	validationError := new(models.ValidationError)
	if c.profileStatisticWeightId == "" {
		validationError.AddError("profileStatisticWeightId", errors.New(models.ArgumentNullOrEmptyError))
	}
	if !utils.IsValidUUID(c.profileStatisticWeightId) {
		validationError.AddError("profileStatisticWeightId", errors.New(models.InvalidFormat))
	}
	if validationError.HasErrors() {
		return validationError
	}
	return nil
}

func (c *ProfileStatisticWeightDeleteByIdCommand) Execute() (interface{}, error) {
	err := c.userStatisticService.DeleteWeightMeasurementById(c.profileStatisticWeightId)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
