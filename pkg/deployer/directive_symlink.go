package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

// Symlink - creates symlink
type Symlink struct {
	From       string
	To         string
	Create     bool
	Delete     bool
	isModified bool
}

const DirectiveNameSymlink = "Symlink"

func (d *Symlink) New() Directive {
	return &Symlink{}
}

func (d *Symlink) Name() string {
	return DirectiveNameSymlink
}

func (d *Symlink) FileName(id int) string {
	return fmt.Sprintf("%d%s%s%s",
		id,
		RunListFileSeparator,
		d.Name(),
		RunListFileExtension,
	)
}

// Hydrate - hydrates attributes from run-list file.data
func (d *Symlink) Hydrate(data []byte) error {
	if data == nil {
		return fmt.Errorf("symlink data can't be nil")
	}

	if err := json.Unmarshal(data, &d); err != nil {
		return fmt.Errorf("symlink data Unmarshal failed with: %w", err)
	}
	return nil
}

// Init - initialize directive with ref to SDK and global attributes
func (d *Symlink) Init(aa *Attributes) error {
	return nil
}

// Execute - creates symlink if one doesn't exist yet
func (d *Symlink) Execute(ctx context.Context, log SimpleLogger, in *DirectiveInput) error {
	if d.Create {
		return d.create(ctx, log, in)
	}
	if d.Delete {
		return d.delete(ctx, log, in)
	}
	return fmt.Errorf("no action specified")
}

func (d *Symlink) IsModified() bool {
	return d.isModified
}

func (d *Symlink) DependsOn() []int {
	// we don't have any do only ifModified dependencies for this directive
	return nil
}

func (d *Symlink) create(ctx context.Context, log SimpleLogger, in *DirectiveInput) error {
	const (
		actionCreate = "create"
		actionSkip   = "skip (already exists)"
	)

	action := actionSkip
	_, err := os.Lstat(d.To)
	if err != nil {
		if os.IsNotExist(err) {
			action = actionCreate
		}
	}

	log.Printf("... ( dryRun=%v ) %s -> %s: %s action=%s",
		in.E.DryRun,
		d.Name(),
		d.From,
		d.To,
		action,
	)

	if in.E.DryRun || action == actionSkip {
		return nil
	}

	if err := os.Symlink(d.From, d.To); err != nil {
		return fmt.Errorf("os.Symlink failed with: %w", err)
	}
	d.isModified = true
	return nil
}

func (d *Symlink) delete(ctx context.Context, log SimpleLogger, in *DirectiveInput) error {
	const (
		actionDelete = "delete"
		actionSkip   = "skip (already deleted)"
	)

	action := actionDelete
	_, err := os.Lstat(d.To)
	if err != nil {
		if os.IsNotExist(err) {
			action = actionSkip
		}
	}

	log.Printf("... ( dryRun=%v ) %s: %s action=%s",
		in.E.DryRun,
		d.Name(),
		d.To,
		action,
	)

	if in.E.DryRun || action == actionSkip {
		return nil
	}

	if err := os.Remove(d.To); err != nil {
		return fmt.Errorf("os.Remove failed with: %w", err)
	}
	d.isModified = true
	return nil
}
