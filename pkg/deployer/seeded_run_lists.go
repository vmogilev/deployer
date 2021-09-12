package deployer

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
			&GenerateFile{
				Path:     "/var/www/html/index.php",
				Owner:    "www-data",
				Group:    "www-data",
				Mode:     0640,
				Template: "hello-world.php",
				Data:     nil,
			},
		},
	},
}
