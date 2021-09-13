package deployer

import (
	"context"
	"os"
	"path/filepath"
)

type RunList struct {
	Name string
	All  []Directive
}

func (l *RunList) RunListName() string {
	return l.Name
}

// BasePath - directory where all run list directives are defined
func (l *RunList) BasePath(configDir string) string {
	return filepath.Join(configDir, DirNameRunList, l.Name)
}

// BuildFromPath - builds run list from os path
func (l *RunList) BuildFromPath(path string, log SimpleLogger) error {
	base := l.BasePath(path)
	rr, err := os.ReadDir(base)
	if err != nil {
		return err
	}

	for _, r := range rr {
		// shouldn't happen but let's be resilient and skip if it does
		if r.IsDir() {
			continue
		}

		log.Printf("... parsing %s", r.Name())

		// parse the run-list file name attributes
		f := &RunListFile{Name: r.Name()}
		if err := f.Parse(); err != nil {
			return err
		}

		// load run-list file
		data, err := os.ReadFile(filepath.Join(base, r.Name()))
		if err != nil {
			return err
		}

		// parse directive and hydrate from run-list file
		d := DirectiveFromName(f.directive)
		if err := d.Hydrate(data); err != nil {
			return err
		}

		l.All = append(l.All, d)
	}

	return nil
}

// ExecuteAll - executes all directives in the run list
func (l *RunList) ExecuteAll(ctx context.Context, log SimpleLogger, globals *Attributes, in *DirectiveInput) error {
	for i, d := range l.All {
		log.Printf("... processing [%d] %s", i, d.Name())

		// initialize directive with globals
		if err := d.Init(globals); err != nil {
			return err
		}

		// we always execute during dryRun or when no deps exist
		var execute bool
		if in.E.DryRun || len(d.DependsOn()) == 0 {
			execute = true
		}

		// check if any deps were modified to see if we need to execute
		if !execute {
			for _, i := range d.DependsOn() {
				if l.All[i].IsModified() {
					execute = true
					break
				}
			}
		}

		if !execute {
			log.Printf("... skipping %s since it's dependency wasn't modified", d)
			continue
		}

		// finally execute directive
		if err := d.Execute(ctx, log, in); err != nil {
			return err
		}
	}

	return nil
}

func (l *RunList) ListAll() []Directive {
	return l.All
}
