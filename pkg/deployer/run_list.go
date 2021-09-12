package deployer

import "path/filepath"

type RunList struct {
	Name string
	All  []Directive
}

func (l *RunList) BasePath(configDir string) string {
	return filepath.Join(configDir, DirNameRunList, l.Name)
}
