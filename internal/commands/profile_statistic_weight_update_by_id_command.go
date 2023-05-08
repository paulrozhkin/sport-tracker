package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"github.com/paulrozhkin/sport-tracker/internal/utils"
)

type ProfileStatisticWeightUpdateByIdCommand struct {
	AuthorizedCommand
	contextModel             *dto.CreateWeightStatisticModel
	context                  *CommandContext
	userStatisticService     *services.UserStatisticService
	profileStatisticWeightId string
}

func NewProfileStatisticWeightUpdateByIdCommand(userStatisticService *services.UserStatisticService) (*ProfileStatisticWeightUpdateByIdCommand, error) {
	command := new(ProfileStatisticWeightUpdateByIdCommand)
	command.contextModel = &dto.CreateWeightStatisticModel{}
	command.context = &CommandContext{CommandParameters: map[string]string{
		"profileStatisticWeightId": "",
	}, CommandContent: command.contextModel}
	command.userStatisticService = userStatisticService
	return command, nil
}

func (c *ProfileStatisticWeightUpdateByIdCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *ProfileStatisticWeightUpdateByIdCommand) Validate() error {
	c.profileStatisticWeightId = c.context.CommandParameters["profileStatisticWeightId"]
	validationError := new(models.ValidationError)
	if c.contextModel.Weight >= services.MaxWeight || c.contextModel.Weight <= services.MinWeight {
		validationError.AddError("weight", errors.New("must be less 3000 or more 1"))
	}
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

func (c *ProfileStatisticWeightUpdateByIdCommand) Execute() (interface{}, error) {
	model := models.UserWeight{User: c.claims.UserId,
		Weight: c.contextModel.Weight,
		Date:   c.contextModel.Date.Time}
	model.Id = c.profileStatisticWeightId
	updatedModel, err := c.userStatisticService.UpdateUserWeightMeasurement(model)
	if err != nil {
		return nil, err
	}
	return mapWeightStatisticToDto(updatedModel), nil
}
