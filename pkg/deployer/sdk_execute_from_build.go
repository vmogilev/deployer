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
// and uses templates embedded in the binary instead of os directory
func (x *SDK) ExecuteFromBuild(ctx context.Context, in *ExecuteInput) error {
	if in.ListOnly {
		for l := range runableAll {
			fmt.Printf("%s\n", l)
		}
		return nil
	}

	if err := x.gen.ParseFromBuild(templatesALL, filepath.Join(DirNameTemplates, "*")); err != nil {
		return fmt.Errorf("call to ParseFromBuild failed with: %w", err)
	}

	l, ok := runableAll[in.RunList]
	if !ok {
		return fmt.Errorf("%s run-list not found in runableAll", in.RunList)
	}

	return l.ExecuteAll(ctx, x.log, x.globals, &DirectiveInput{
		G:       x.gen,
		E:       in,
		Verbose: x.in.Verbose,
	})
}
