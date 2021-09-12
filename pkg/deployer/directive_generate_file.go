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
	Path     string
	Owner    string
	Group    string
	Mode     os.FileMode
	Template string
	Data     map[string]interface{}
	x        *SDK
}

const DirectiveNameGenerateFile = "GenerateFile"

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
func (d *GenerateFile) Init(x *SDK, aa *Attributes) error {
	d.x = x
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
func (d *GenerateFile) Execute(ctx context.Context, in *ExecuteInput) error {
	var buf bytes.Buffer
	if err := d.x.gen.Execute(&buf, d.Template, d.Data); err != nil {
		return fmt.Errorf("GenerateFile failed to execute %s template with: %w", d.Template, err)
	}

	const (
		actionCreate = "create"
		actionUpdate = "update"
	)

	action := actionUpdate
	_, err := os.Stat(d.Path)
	if err != nil {
		if os.IsNotExist(err) {
			action = actionCreate
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

	d.x.log.Printf("... ( dryRun=%v ) %s: %s action=%s",
		in.DryRun,
		d.Name(),
		d.Path,
		action,
	)
	d.x.log.Printf("... ( dryRun=%v ) %s: %s chmod=%s",
		in.DryRun,
		d.Name(),
		d.Path,
		d.Mode,
	)
	d.x.log.Printf("... ( dryRun=%v ) %s: %s chown=%s:%s (%d:%d)",
		in.DryRun,
		d.Name(),
		d.Path,
		d.Owner,
		d.Group,
		uid,
		gid,
	)
	if d.x.in.Verbose {
		fmt.Println(buf.String())
	}
	if in.DryRun {
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
	return nil
}
