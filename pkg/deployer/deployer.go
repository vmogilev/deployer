package deployer

import "context"

type SimpleLogger interface {
	Printf(string, ...interface{})
}
type SDK struct {
	in  *SDKInput
	log SimpleLogger
}

type SDKInput struct {
	Verbose bool
}

func NewSDK(log SimpleLogger, in *SDKInput) *SDK {
	return &SDK{
		in:  in, // input pass through to minimize refactoring
		log: log,
	}
}

type ExecuteInput struct {
	DryRun bool
	Force  bool
}

func (x *SDK) Execute(ctx context.Context, in *ExecuteInput) error {
	return nil
}
