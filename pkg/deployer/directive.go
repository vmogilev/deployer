package deployer

import (
	"context"
	"fmt"
)

type Directive interface {
	Name() string
	FileName(in int) string
	Hydrate(data []byte) error
	Init(x *SDK, aa *Attributes) error
	Execute(ctx context.Context, in *ExecuteInput) error
}

type unimplementedDirective struct {
	name string
}

func (d *unimplementedDirective) Name() string {
	return fmt.Sprintf("%s directive not implemented", d.name)
}

func (d *unimplementedDirective) FileName(in int) string {
	return fmt.Sprintf("%s directive not implemented", d.name)
}

func (d *unimplementedDirective) Hydrate(data []byte) error {
	return fmt.Errorf(d.Name())
}

func (d *unimplementedDirective) Init(x *SDK, aa *Attributes) error {
	return fmt.Errorf(d.Name())
}

func (d *unimplementedDirective) Execute(ctx context.Context, in *ExecuteInput) error {
	return fmt.Errorf(d.Name())
}

// DirectiveFromName - constructs new directive by its name
func DirectiveFromName(name string) Directive {
	switch name {
	case DirectiveNameGenerateFile:
		return &GenerateFile{}
	default:
		return &unimplementedDirective{name: name}
	}
}
