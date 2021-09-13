package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// OsPackage - manages OsPackage
type OsPackage struct {
	PkgName    string
	Install    bool
	Remove     bool
	isModified bool
}

const DirectiveNameOsPackage = "OsPackage"

func (d *OsPackage) New() Directive {
	return &OsPackage{}
}

func (d *OsPackage) Name() string {
	return DirectiveNameOsPackage
}

func (d *OsPackage) FileName(id int) string {
	return fmt.Sprintf("%d%s%s%s",
		id,
		RunListFileSeparator,
		d.Name(),
		RunListFileExtension,
	)
}

// Hydrate - hydrates attributes from run-list file.data
func (d *OsPackage) Hydrate(data []byte) error {
	if data == nil {
		return fmt.Errorf("OsPackage data can't be nil")
	}

	if err := json.Unmarshal(data, &d); err != nil {
		return fmt.Errorf("OsPackage data Unmarshal failed with: %w", err)
	}
	return nil
}

// Init - initialize directive with ref to SDK and global attributes
func (d *OsPackage) Init(aa *Attributes) error {
	return nil
}

// Execute - creates OsPackage if one doesn't exist yet
func (d *OsPackage) Execute(ctx context.Context, log SimpleLogger, in *DirectiveInput) error {
	if d.Install {
		return d.install(ctx, log, in)
	}
	if d.Remove {
		return d.remove(ctx, log, in)
	}
	return fmt.Errorf("no action specified")
}

func (d *OsPackage) IsModified() bool {
	return d.isModified
}

func (d *OsPackage) DependsOn() []int {
	// we don't have any do only ifModified dependencies for this directive
	return nil
}

func (d *OsPackage) install(ctx context.Context, log SimpleLogger, in *DirectiveInput) error {
	const (
		actionInstall = "install"
		actionSkip    = "skip (already installed)"
	)

	action := actionSkip
	installed, err := d.IsInstalled(log, in.Verbose)
	if err != nil {
		return err
	}
	if !installed {
		action = actionInstall
	}

	log.Printf("... ( dryRun=%v ) %s %s action=%s",
		in.E.DryRun,
		d.Name(),
		d.PkgName,
		action,
	)

	if in.E.DryRun || action == actionSkip {
		return nil
	}

	return d.aptGet(log, action, in.Verbose)
}

func (d *OsPackage) remove(ctx context.Context, log SimpleLogger, in *DirectiveInput) error {
	const (
		actionRemove = "remove"
		actionSkip   = "skip (already removed)"
	)

	action := actionSkip
	installed, err := d.IsInstalled(log, in.Verbose)
	if err != nil {
		return err
	}
	if installed {
		action = actionRemove
	}

	log.Printf("... ( dryRun=%v ) %s %s action=%s",
		in.E.DryRun,
		d.Name(),
		d.PkgName,
		action,
	)

	if in.E.DryRun || action == actionSkip {
		return nil
	}

	return d.aptGet(log, action, in.Verbose)
}

func (d *OsPackage) IsInstalled(log SimpleLogger, verbose bool) (bool, error) {
	cmd := exec.Command("dpkg", "-l", d.PkgName)
	out, err := cmd.CombinedOutput()
	if verbose {
		log.Printf(string(out))
	}
	if err == nil {
		return true, nil
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		if strings.Contains(string(out), "no packages found matching") {
			return false, nil
		}
		log.Printf("ERROR: %s", string(out))
		return false, fmt.Errorf("dpkg -l %s finished with non-zero code: %w",
			d.PkgName,
			exitErr,
		)
	}
	log.Printf("ERROR: %s", string(out))
	return false, fmt.Errorf("failed to run dpkg -l %s: %w",
		d.PkgName,
		err,
	)
}

func (d *OsPackage) aptGet(log SimpleLogger, arg string, verbose bool) error {
	cmd := exec.Command("apt-get", arg, "-y", d.PkgName)
	out, err := cmd.CombinedOutput()
	if verbose {
		log.Printf(string(out))
	}
	if err == nil {
		d.isModified = true
		return nil
	}

	log.Printf("ERROR: %s", string(out))
	if exitErr, ok := err.(*exec.ExitError); ok {
		return fmt.Errorf("apt-get %s %s finished with non-zero code: %w",
			arg,
			d.PkgName,
			exitErr,
		)
	}
	return fmt.Errorf("failed to run apt-get %s %s: %w",
		arg,
		d.PkgName,
		err,
	)
}
