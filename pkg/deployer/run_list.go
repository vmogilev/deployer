package deployer

type RunList struct {
	Name string
	All  []Directive
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
				OwnerID:  33, // www-data
				GroupID:  33, // www-data
				Mode:     0640,
				Template: "hello-world.com",
				Data:     nil,
			},
		},
	},
}
