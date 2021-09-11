package deployer

import (
	"context"
	"fmt"
)

type Directive interface {
	Name() string
	Init(data []byte, aa *Attributes) error
	Execute(ctx context.Context, in *ExecuteInput) error
}

type unimplementedDirective struct {
	name string
}

func (d *unimplementedDirective) Name() string {
	return fmt.Sprintf("%s directive not implemented", d.name)
}

func (d *unimplementedDirective) Init(data []byte, aa *Attributes) error {
	return fmt.Errorf(d.Name())
}

func (d *unimplementedDirective) Execute(ctx context.Context, in *ExecuteInput) error {
	return fmt.Errorf(d.Name())
}

// DirectiveFromName - constructs new directive by its name
func (x *SDK) DirectiveFromName(name string) Directive {
	switch name {
	case DirectiveNameGenerateFile:
		return &GenerateFile{x: x}
	default:
		return &unimplementedDirective{name: name}
	}
}
