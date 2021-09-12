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
	base := l.BasePath(x.in.Vars.DeployerConfigDir)
	rr, err := os.ReadDir(base)
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
		data, err := os.ReadFile(filepath.Join(base, r.Name()))
		if err != nil {
			return err
		}

		// parse directive and hydrate from run-list file
		d := DirectiveFromName(f.directive)
		if err := d.Hydrate(data); err != nil {
			return err
		}

		// initialize directive with globals
		if err := d.Init(x, gg); err != nil {
			return err
		}

		// finally execute directive
		if err := d.Execute(ctx, in); err != nil {
			return err
		}
	}

	return nil
}
