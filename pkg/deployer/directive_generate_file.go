package deployer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type GenerateFile struct {
	Path     string
	OwnerID  int
	GroupID  int
	Mode     os.FileMode
	Template string
	Data     map[string]interface{}
	x        *SDK
}

const DirectiveNameGenerateFile = "GenerateFile"

func (d *GenerateFile) Name() string {
	return DirectiveNameGenerateFile
}

// Init - hydrates attributes from run-list file.data
func (d *GenerateFile) Init(data []byte, aa *Attributes) error {
	if data == nil {
		return fmt.Errorf("GenerateFile data can't be nil")
	}

	if err := json.Unmarshal(data, &d); err != nil {
		return fmt.Errorf("GenerateFile data Unmarshal failed with: %w", err)
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

	d.x.log.Printf("... ( dryRun=%v ) %s: action=%s path=%s ",
		in.DryRun,
		d.Name(),
		action,
		d.Path,
	)
	d.x.log.Printf("... ( dryRun=%v ) %s: chmod=%s path=%s ",
		in.DryRun,
		d.Name(),
		d.Mode,
		d.Path,
	)
	d.x.log.Printf("... ( dryRun=%v ) %s: chown=%d:%d path=%s ",
		in.DryRun,
		d.Name(),
		d.OwnerID,
		d.GroupID,
		d.Path,
	)
	if in.DryRun {
		return nil
	}

	if err := os.WriteFile(d.Path, buf.Bytes(), d.Mode); err != nil {
		return fmt.Errorf("writeFile failed with: %w", err)
	}
	if err := os.Chmod(d.Path, d.Mode); err != nil {
		return fmt.Errorf("chmod failed with: %w", err)
	}
	if err := os.Chown(d.Path, d.OwnerID, d.GroupID); err != nil {
		return fmt.Errorf("chown failed with: %w", err)
	}
	return nil
}
