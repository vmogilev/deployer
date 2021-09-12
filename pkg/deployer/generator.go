package deployer

import (
	"io"
	"io/fs"
	"path/filepath"
	"text/template"
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

// ParseFromBuild - parses all templates compiled at build-time in go:embed fs/*
func (g *Generator) ParseFromBuild(fs fs.FS, pat string) error {
	var err error
	g.all[BaseTemplates], err = template.New(pat).ParseFS(fs, pat)
	return err
}

// ParseFromFS - parses all templates from local file system
func (g *Generator) ParseFromFS(subDir string) error {
	var err error
	g.all[BaseTemplates], err = template.New(subDir).ParseGlob(filepath.Join(g.rootDir, subDir, "*"))
	return err
}

// Execute - executes previously parsed template by its name, pushing result to wr
func (g *Generator) Execute(wr io.Writer, name string, data interface{}) error {
	return g.all[BaseTemplates].ExecuteTemplate(wr, name, data)
}
