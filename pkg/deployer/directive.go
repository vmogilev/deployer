package deployer

import (
	"context"
	"fmt"
)

type DirectiveInput struct {
	G       FileGenerator
	E       *ExecuteInput
	Verbose bool
}

type Directive interface {
	New() Directive
	Name() string
	FileName(in int) string
	Hydrate(data []byte) error
	Init(aa *Attributes) error
	Execute(ctx context.Context, log SimpleLogger, in *DirectiveInput) error

	// IsModified - indicates to downstream dep that
	// something was modified at OS, and it should take action
	IsModified() bool

	// DependsOn - list of IDs of this run list directives
	// this directive depends on being in modified state
	// if none of these directives are in modified state
	// this directive `Execute()` phase is skipped ...
	DependsOn() []int
}

type unimplementedDirective struct {
	name string
}

func (d *unimplementedDirective) New() Directive {
	return &unimplementedDirective{}
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

func (d *unimplementedDirective) Init(aa *Attributes) error {
	return fmt.Errorf(d.Name())
}

func (d *unimplementedDirective) Execute(ctx context.Context, log SimpleLogger, in *DirectiveInput) error {
	return fmt.Errorf(d.Name())
}

func (d *unimplementedDirective) IsModified() bool {
	return false
}

func (d unimplementedDirective) DependsOn() []int {
	return nil
}

var directivesAll = map[string]Directive{}

// RegisterDirective - registers directive in directivesAll
func RegisterDirective(d Directive) {
	directivesAll[d.Name()] = d
}

// DirectiveFromName - constructs new directive by its name
func DirectiveFromName(name string) Directive {
	if d, ok := directivesAll[name]; ok {
		return d.New()
	}
	return &unimplementedDirective{name: name}
}
