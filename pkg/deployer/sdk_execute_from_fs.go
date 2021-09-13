package deployer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// ExecuteFromFS - executes All directives from file system run-list in the
// order of their ID parsed from the file name {ID}_{DIRECTIVE}.json
// The {ID}_{DIRECTIVE}.json contains the required Attributes and
// references to file templates if any ...
func (x *SDK) ExecuteFromFS(ctx context.Context, in *ExecuteInput) error {
	if in.ListOnly {
		p := filepath.Join(x.in.Vars.DeployerConfigDir, DirNameRunList)
		ll, err := os.ReadDir(p)
		if err != nil {
			return err
		}
		for _, l := range ll {
			if l.IsDir() {
				fmt.Printf("%s\n", l)
			}
		}
		return nil
	}

	if err := x.gen.ParseFromFS(DirNameTemplates); err != nil {
		return err
	}

	l := &RunList{Name: in.RunList}
	if err := l.BuildFromPath(x.in.Vars.DeployerConfigDir, x.log); err != nil {
		return err
	}

	return l.ExecuteAll(ctx, x.log, x.globals, &DirectiveInput{
		G:       x.gen,
		E:       in,
		Verbose: x.in.Verbose,
	})
}
