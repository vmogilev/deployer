package deployer

import (
	"context"
	"embed"
	"io"
	"os"
	"path/filepath"

	"github.com/vmogilev/deployer/internal/env"
)

const (
	DirNameRunList   = "run"
	DirNameTemplates = "tmpl"
	DirNameCache     = ".cache"
	BaseTemplates    = "base"
)

type SimpleLogger interface {
	Printf(string, ...interface{})
}

type FileGenerator interface {
	Execute(wr io.Writer, name string, data interface{}) error
}

type SDK struct {
	in  *SDKInput
	log SimpleLogger
	gen FileGenerator
}

type SDKInput struct {
	Vars    *env.Vars
	Verbose bool
}

//go:embed templates/*
var templatesALL embed.FS

func NewSDK(log SimpleLogger, in *SDKInput) (*SDK, error) {
	gen := NewGenerator(in.Vars.DeployerConfigDir)
	if err := gen.Parse(templatesALL, BaseTemplates); err != nil {
		return nil, err
	}

	return &SDK{
		in:  in, // input pass through to minimize refactoring
		log: log,
		gen: gen,
	}, nil
}

type ExecuteInput struct {
	RunList string
	DryRun  bool
	Force   bool
}

// Execute - executes All directives from a given run-list in the
// order of their ID parsed from the file name {ID}_{DIRECTIVE}.json
// The {ID}_{DIRECTIVE}.json contains the required Attributes and
// references to file templates if any ...
func (x *SDK) Execute(ctx context.Context, in *ExecuteInput) error {
	p := filepath.Join(x.in.Vars.DeployerConfigDir, DirNameRunList, in.RunList)
	rr, err := os.ReadDir(p)
	if err != nil {
		return err
	}

	gg, err := x.Globals()
	if err != nil {
		return err
	}

	for _, r := range rr {
		// shouldn't happen but let's be resilient and skip if it does
		if r.IsDir() {
			continue
		}

		x.log.Printf("... processing %s", r.Name())

		// parse the run-list file name attributes
		f := &RunListFile{Name: r.Name()}
		if err := f.Parse(); err != nil {
			return err
		}

		// load run-list file
		data, err := os.ReadFile(r.Name())
		if err != nil {
			return err
		}

		// parse/init directive
		d := x.DirectiveFromName(f.directive)
		if err := d.Init(data, gg); err != nil {
			return err
		}

		// finally execute directive
		if err := d.Execute(ctx, in); err != nil {
			return err
		}
	}

	return nil
}
