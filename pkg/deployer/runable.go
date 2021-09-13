package deployer

import "context"

type Runable interface {
	RunListName() string
	BasePath(configDir string) string
	BuildFromPath(path string, log SimpleLogger) error
	ExecuteAll(ctx context.Context, log SimpleLogger, globals *Attributes, in *DirectiveInput) error
	ListAll() []Directive
}

var runableAll = map[string]Runable{}

// RegisterRunable - registers runable run-list
func RegisterRunable(r Runable) {
	runableAll[r.RunListName()] = r
}
