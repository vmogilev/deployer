package deployer

import "path/filepath"

type RunList struct {
	Name string
	All  []Directive
}

func (l *RunList) BasePath(configDir string) string {
	return filepath.Join(configDir, DirNameRunList, l.Name)
}

const (
	RunListHelloWorld = "hello-world"
)

var seededRunLists = map[string]*RunList{
	RunListHelloWorld: {
		Name: RunListHelloWorld,
		All: []Directive{
			&GenerateFile{
				Path:     "/etc/nginx/sites-available/hello-world.com",
				Owner:    "www-data",
				Group:    "www-data",
				Mode:     0640,
				Template: "nginx-site",
				Data:     nil,
			},
		},
	},
}
