package deployer

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/vmogilev/deployer/internal/env"
)

const (
	DirNameRunList   = "run"
	DirNameTemplates = "templates"
	DirNameCache     = ".cache"
	BaseTemplates    = "base"
)

type SimpleLogger interface {
	Printf(string, ...interface{})
}

type FileGenerator interface {
	ParseFromBuild(fs fs.FS, name string) error
	ParseFromFS(subDir string) error
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

func NewSDK(log SimpleLogger, in *SDKInput) *SDK {
	return &SDK{
		in:  in, // input pass through to minimize refactoring
		log: log,
		gen: NewGenerator(in.Vars.DeployerConfigDir),
	}
}

type ExecuteInput struct {
	RunList  string
	DryRun   bool
	Force    bool
	ListOnly bool
}

func (i *ExecuteInput) Validate() error {
	if i.ListOnly {
		return nil
	}
	if i.RunList == "" {
		return fmt.Errorf("runList is a required parameter")
	}
	return nil
}
