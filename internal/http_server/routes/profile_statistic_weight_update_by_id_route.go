package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type ProfileStatisticWeightUpdateByIdRoute struct {
	userStatisticService *services.UserStatisticService
}

func NewProfileStatisticWeightUpdateByIdRoute(userStatisticService *services.UserStatisticService) *ProfileStatisticWeightUpdateByIdRoute {
	return &ProfileStatisticWeightUpdateByIdRoute{userStatisticService: userStatisticService}
}

func (*ProfileStatisticWeightUpdateByIdRoute) Method() string {
	return http.MethodPut
}

func (*ProfileStatisticWeightUpdateByIdRoute) Pattern() string {
	return "/profile/statistic/weight/{profileStatisticWeightId}"
}

func (a *ProfileStatisticWeightUpdateByIdRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewProfileStatisticWeightUpdateByIdCommand(a.userStatisticService)
}
