package commands

import (
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
)

type ProfileStatisticWeightCreateCommand struct {
	AuthorizedCommand
	contextModel         *dto.CreateWeightStatisticModel
	context              *CommandContext
	userStatisticService *services.UserStatisticService
}

func NewProfileStatisticWeightCreateCommand(userStatisticService *services.UserStatisticService) (*ProfileStatisticWeightCreateCommand, error) {
	contextModel := &dto.CreateWeightStatisticModel{}
	context := &CommandContext{CommandContent: contextModel}
	return &ProfileStatisticWeightCreateCommand{context: context, userStatisticService: userStatisticService, contextModel: contextModel}, nil
}

func (c *ProfileStatisticWeightCreateCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *ProfileStatisticWeightCreateCommand) Validate() error {
	validationError := new(models.ValidationError)
	if c.contextModel.Weight >= services.MaxWeight || c.contextModel.Weight <= services.MinWeight {
		validationError.AddError("weight", errors.New("must be less 3000 or more 1"))
	}
	if validationError.HasErrors() {
		return validationError
	}
	return nil
}

func (c *ProfileStatisticWeightCreateCommand) Execute() (interface{}, error) {
	createModel := models.UserWeight{User: c.claims.UserId,
		Weight: c.contextModel.Weight,
		Date:   c.contextModel.Date.Time}
	createdModel, err := c.userStatisticService.CreateUserWeightMeasurement(createModel)
	if err != nil {
		return nil, err
	}
	return mapWeightStatisticToDto(createdModel), nil
}

func mapWeightStatisticToDto(weightMeasurement *models.UserWeight) *dto.WeightStatisticModel {
	return &dto.WeightStatisticModel{
		Weight: weightMeasurement.Weight,
		Id:     weightMeasurement.Id,
		Date:   dto.JsonDate{Time: weightMeasurement.Date},
	}
}
