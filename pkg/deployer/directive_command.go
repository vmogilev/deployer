package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

// Command - creates Command
type Command struct {
	Run          string
	Dependencies []int
}

const DirectiveNameCommand = "Command"

func (d *Command) New() Directive {
	return &Command{}
}

func (d *Command) Name() string {
	return DirectiveNameCommand
}

func (d *Command) FileName(id int) string {
	return fmt.Sprintf("%d%s%s%s",
		id,
		RunListFileSeparator,
		d.Name(),
		RunListFileExtension,
	)
}

// Hydrate - hydrates attributes from run-list file.data
func (d *Command) Hydrate(data []byte) error {
	if data == nil {
		return fmt.Errorf("command data can't be nil")
	}

	if err := json.Unmarshal(data, &d); err != nil {
		return fmt.Errorf("command data Unmarshal failed with: %w", err)
	}
	return nil
}

// Init - initialize directive with ref to SDK and global attributes
func (d *Command) Init(aa *Attributes) error {
	return nil
}

// Execute - calls exec.Command
func (d *Command) Execute(ctx context.Context, log SimpleLogger, in *DirectiveInput) error {
	log.Printf("... ( dryRun=%v ) %s: %s (depends on %v)",
		in.E.DryRun,
		d.Name(),
		d.Run,
		d.Dependencies,
	)
	if in.E.DryRun {
		return nil
	}

	cmd := exec.Command("bash", "-c", d.Run)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("ERROR: %s", string(out))
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("running %q finished with non-zero code: %w",
				d.Run,
				exitErr,
			)
		}
		return fmt.Errorf("failed to run bash -c %s: %w",
			d.Run,
			err,
		)
	}

	if in.Verbose {
		log.Printf(string(out))
	}
	return nil
}

func (d *Command) IsModified() bool {
	return true
}

func (d *Command) DependsOn() []int {
	return d.Dependencies
}

func (d *Command) String() string {
	return fmt.Sprintf("%s(%s)", d.Name(), d.Run)
}
