package commands

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands/dto"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/services"
)

type ProfileStatisticGetCommand struct {
	AuthorizedCommand
	context              *CommandContext
	userStatisticService *services.UserStatisticService
}

func NewProfileStatisticGetCommand(userStatisticService *services.UserStatisticService) (*ProfileStatisticGetCommand, error) {
	context := &CommandContext{}
	return &ProfileStatisticGetCommand{context: context, userStatisticService: userStatisticService}, nil
}

func (c *ProfileStatisticGetCommand) GetCommandContext() *CommandContext {
	return c.context
}

func (c *ProfileStatisticGetCommand) Validate() error {
	return nil
}

func (c *ProfileStatisticGetCommand) Execute() (interface{}, error) {
	statistic, err := c.userStatisticService.GetGeneralStatisticForUser(c.claims.UserId)
	if err != nil {
		return nil, err
	}
	return mapUserStatisticToDto(statistic), nil
}

func mapUserStatisticToDto(statistic *models.UserStatistic) *dto.GeneralStatisticModel {
	result := new(dto.GeneralStatisticModel)
	result.WeightStatistic = make([]*dto.WeightStatisticModel, 0)
	for _, item := range statistic.Weight {
		result.WeightStatistic = append(result.WeightStatistic, mapWeightStatisticToDto(item))
	}
	result.WorkoutsPerMonth = statistic.WorkoutsPerMonth
	result.WorkoutsPerYear = statistic.WorkoutsPerYear
	return result
}
