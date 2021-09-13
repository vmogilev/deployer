package deployer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strconv"
)

// GenerateFile - generates files from template
// TODO: Mode serializes as decimal and it'd be nice to retain octal format
type GenerateFile struct {
	Path         string
	Owner        string
	Group        string
	Mode         os.FileMode
	Template     string
	Data         map[string]interface{}
	isModified   bool
	Dependencies []int
}

const DirectiveNameGenerateFile = "GenerateFile"

func (d *GenerateFile) New() Directive {
	return &GenerateFile{}
}

func (d *GenerateFile) Name() string {
	return DirectiveNameGenerateFile
}

func (d *GenerateFile) FileName(id int) string {
	return fmt.Sprintf("%d%s%s%s",
		id,
		RunListFileSeparator,
		d.Name(),
		RunListFileExtension,
	)
}

// Hydrate - hydrates attributes from run-list file.data
func (d *GenerateFile) Hydrate(data []byte) error {
	if data == nil {
		return fmt.Errorf("GenerateFile data can't be nil")
	}

	if err := json.Unmarshal(data, &d); err != nil {
		return fmt.Errorf("GenerateFile data Unmarshal failed with: %w", err)
	}
	return nil
}

// Init - appends attributes to data
func (d *GenerateFile) Init(aa *Attributes) error {
	if d.Data == nil {
		d.Data = make(map[string]interface{})
	}
	// merge local data with global attributes, having local carry higher weight
	for k, v := range aa.All {
		if _, ok := d.Data[k]; ok {
			continue
		}
		d.Data[k] = v
	}
	return nil
}

// Execute - applies changes to managed file
func (d *GenerateFile) Execute(ctx context.Context, log SimpleLogger, in *DirectiveInput) error {
	var buf bytes.Buffer
	if err := in.G.Execute(&buf, d.Template, d.Data); err != nil {
		return fmt.Errorf("GenerateFile failed to execute %s template with: %w", d.Template, err)
	}

	const (
		actionCreate = "create"
		actionUpdate = "update"
		actionSkip   = "skip (file unchanged)"
	)

	action := actionUpdate
	_, err := os.Stat(d.Path)
	if err != nil {
		if os.IsNotExist(err) {
			action = actionCreate
		}
	}

	if action == actionUpdate {
		oldData, err := os.ReadFile(d.Path)
		if err != nil {
			return err
		}
		if res := bytes.Compare(oldData, buf.Bytes()); res == 0 {
			action = actionSkip
		}
	}

	u, err := user.Lookup(d.Owner)
	if err != nil {
		return err
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return err
	}

	g, err := user.LookupGroup(d.Group)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return err
	}

	log.Printf("... ( dryRun=%v ) %s: %s action=%s",
		in.E.DryRun,
		d.Name(),
		d.Path,
		action,
	)
	log.Printf("... ( dryRun=%v ) %s: %s chmod=%s",
		in.E.DryRun,
		d.Name(),
		d.Path,
		d.Mode,
	)
	log.Printf("... ( dryRun=%v ) %s: %s chown=%s:%s (%d:%d)",
		in.E.DryRun,
		d.Name(),
		d.Path,
		d.Owner,
		d.Group,
		uid,
		gid,
	)
	if in.Verbose {
		fmt.Println(buf.String())
	}
	if in.E.DryRun || action == actionSkip {
		// TODO: on actionSkip, possibly look into still doing Chmod and Chown
		return nil
	}

	if err := os.WriteFile(d.Path, buf.Bytes(), d.Mode); err != nil {
		return fmt.Errorf("writeFile failed with: %w", err)
	}
	if err := os.Chmod(d.Path, d.Mode); err != nil {
		return fmt.Errorf("chmod failed with: %w", err)
	}
	if err := os.Chown(d.Path, uid, gid); err != nil {
		return fmt.Errorf("chown failed with: %w", err)
	}
	d.isModified = true
	return nil
}

func (d *GenerateFile) IsModified() bool {
	return d.isModified
}

func (d *GenerateFile) DependsOn() []int {
	return d.Dependencies
}
