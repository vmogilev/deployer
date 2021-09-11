package deployer

import (
	"html/template"
	"io"
	"io/fs"
)

// Generator maintains a list of all known file templates
type Generator struct {
	rootDir string
	all     map[string]*template.Template
}

// NewGenerator - bootstraps new Generator
func NewGenerator(rootDir string) *Generator {
	return &Generator{
		rootDir: rootDir,
		all:     make(map[string]*template.Template),
	}
}

// Parse - parses all templates found in fs/*
func (t *Generator) Parse(fs fs.FS, name string) error {
	var err error
	//t.all[subDir], err = template.New(subDir).ParseGlob(filepath.Join(t.rootDir, subDir, "*"))
	t.all[name], err = template.New(name).ParseFS(fs, "*")
	return err
}

// Execute - executes previously parsed template by its name, pushing result to wr
func (t *Generator) Execute(wr io.Writer, name string, data interface{}) error {
	return t.all[DirNameTemplates].ExecuteTemplate(wr, name, data)
}
