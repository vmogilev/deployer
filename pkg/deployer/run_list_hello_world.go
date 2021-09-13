package deployer

const (
	RunListHelloWorldName = "hello-world"
)

var RunListHelloWorld = &RunList{
	Name: RunListHelloWorldName,
	All: []Directive{
		&Command{
			Run: "apt-get update",
		},
		&GenerateFile{
			Path:     "/etc/nginx/sites-available/hello-world.com",
			Owner:    "www-data",
			Group:    "www-data",
			Mode:     0640,
			Template: "nginx-site",
			Data:     nil,
		},
		&Symlink{
			Create: true,
			From:   "/etc/nginx/sites-available/hello-world.com",
			To:     "/etc/nginx/sites-enabled/hello-world.com",
		},
		&Symlink{
			Delete: true,
			To:     "/etc/nginx/sites-enabled/default",
		},
		&GenerateFile{
			Path:     "/var/www/html/index.php",
			Owner:    "www-data",
			Group:    "www-data",
			Mode:     0640,
			Template: "hello-world.php",
			Data:     nil,
		},
		&Command{
			Run: "nginx -t",
			// we don't care testing nginx unless /etc/nginx/sites-available/hello-world.com was modified
			Dependencies: []int{1},
		},
		&Command{
			Run: "systemctl reload nginx",
			// no need to reload nginx unless /etc/nginx/sites-available/hello-world.com was modified
			Dependencies: []int{1},
		},
	},
}
