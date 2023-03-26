package commands

type CommandContext struct {
	CommandContent    interface{}
	CommandParameters map[string]interface{}
}

type ICommand interface {
	GetCommandContext() *CommandContext
	Validate() error
	Execute() (interface{}, error)
	RequireAuthorization() bool
}
