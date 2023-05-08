package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type ProfileStatisticWeightCreateRoute struct {
	userStatisticService *services.UserStatisticService
}

func NewProfileStatisticWeightCreateRoute(userStatisticService *services.UserStatisticService) *ProfileStatisticWeightCreateRoute {
	return &ProfileStatisticWeightCreateRoute{userStatisticService: userStatisticService}
}

func (*ProfileStatisticWeightCreateRoute) Method() string {
	return http.MethodPost
}

func (*ProfileStatisticWeightCreateRoute) Pattern() string {
	return "/profile/statistic/weight"
}

func (a *ProfileStatisticWeightCreateRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewProfileStatisticWeightCreateCommand(a.userStatisticService)
}
