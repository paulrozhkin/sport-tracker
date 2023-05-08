package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type ProfileStatisticWeightDeleteByIdRoute struct {
	userStatisticService *services.UserStatisticService
}

func NewProfileStatisticWeightDeleteByIdRoute(userStatisticService *services.UserStatisticService) *ProfileStatisticWeightDeleteByIdRoute {
	return &ProfileStatisticWeightDeleteByIdRoute{userStatisticService: userStatisticService}
}

func (*ProfileStatisticWeightDeleteByIdRoute) Method() string {
	return http.MethodDelete
}

func (*ProfileStatisticWeightDeleteByIdRoute) Pattern() string {
	return "/profile/statistic/weight/{profileStatisticWeightId}"
}

func (a *ProfileStatisticWeightDeleteByIdRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewProfileStatisticWeightDeleteByIdCommand(a.userStatisticService)
}
