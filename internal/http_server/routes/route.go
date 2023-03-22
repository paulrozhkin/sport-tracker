package routes

import (
	"github.com/paulrozhkin/sport-tracker/internal/commands"
	"go.uber.org/fx"
)

type Route interface {
	Method() string
	Pattern() string
	NewRouteExecutor() (commands.ICommand, error)
}

// AsRoute annotates the given constructor to state that
// it provides a route to the "routes" group.
func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}
