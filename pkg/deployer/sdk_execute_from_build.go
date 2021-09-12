package deployer

import (
	"context"
	"embed"
	"fmt"
	"path/filepath"
)

//go:embed templates
var templatesALL embed.FS

// ExecuteFromBuild - executes All directives from seededRunLists
func (x *SDK) ExecuteFromBuild(ctx context.Context, in *ExecuteInput) error {
	if in.ListOnly {
		for l := range seededRunLists {
			fmt.Printf("%s\n", l)
		}
		return nil
	}

	if err := x.gen.ParseFromBuild(templatesALL, filepath.Join(DirNameTemplates, "*")); err != nil {
		return fmt.Errorf("call to ParseFromBuild failed with: %w", err)
	}

	gg, err := x.Globals()
	if err != nil {
		return err
	}

	l, ok := seededRunLists[in.RunList]
	if !ok {
		return fmt.Errorf("%s run-list not found in seededRunLists", in.RunList)
	}

	for _, d := range l.All {
		x.log.Printf("... processing %s", d.Name())

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
