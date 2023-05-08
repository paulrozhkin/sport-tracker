package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"github.com/paulrozhkin/sport-tracker/internal/services"
	"net/http"
)

type ProfileStatisticGetRoute struct {
	userStatisticService *services.UserStatisticService
}

func NewProfileStatisticGetRoute(userStatisticService *services.UserStatisticService) *ProfileStatisticGetRoute {
	return &ProfileStatisticGetRoute{userStatisticService: userStatisticService}
}

func (*ProfileStatisticGetRoute) Method() string {
	return http.MethodGet
}

func (*ProfileStatisticGetRoute) Pattern() string {
	return "/profile/statistic/"
}

func (a *ProfileStatisticGetRoute) NewRouteExecutor() (commands.ICommand, error) {
	return commands.NewProfileStatisticGetCommand(a.userStatisticService)
}
