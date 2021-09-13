package deployer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (x *SDK) Export(list string) error {
	// export run-list
	l, ok := runableAll[list]
	if !ok {
		return fmt.Errorf("list %s not found", list)
	}

	base := l.BasePath(x.in.Vars.DeployerConfigDir)
	if err := os.MkdirAll(base, 0700); err != nil {
		return err
	}

	for id, d := range l.ListAll() {
		name := d.FileName(id)
		path := filepath.Join(base, name)
		x.log.Printf("... exporting %s", path)
		data, err := json.MarshalIndent(d, "", "    ")
		if err != nil {
			return fmt.Errorf("error Marshaling: %w", err)
		}
		if err := os.WriteFile(path, data, 0640); err != nil {
			return fmt.Errorf("writeFile failed with: %w", err)
		}
	}

	// export templates
	tt, err := templatesALL.ReadDir(DirNameTemplates)
	if err != nil {
		return fmt.Errorf("can't read embdeded %s dir", DirNameTemplates)
	}
	for _, t := range tt {
		if t.IsDir() {
			continue
		}
		data, err := templatesALL.ReadFile(filepath.Join(DirNameTemplates, t.Name()))
		if err != nil {
			return err
		}
		path := filepath.Join(x.in.Vars.DeployerConfigDir, DirNameTemplates, t.Name())
		x.log.Printf("... exporting %s", path)
		if err := os.WriteFile(path, data, 0640); err != nil {
			return fmt.Errorf("writeFile failed with: %w", err)
		}
	}

	return nil
}
