# Deployer
Deployer is configuration management tool / framework.

Out of the box, it provides the following functionality:
1. create files with template based contents and metadata (owner, group, mode)
2. install and remove Debian packages
3. restart a service when relevant files or packages are updated
4. idempotency
5. demo run-list that installs `hello-world` application (`nginx`->`php-fpm`)
6. embedded run-list export
7. run-list import via Deployer SDK `deployer.RegisterDirective()` and `deployer.RegisterRunable()` (see [Architecture](#architecture))
8. run-list definition from `$DEPLOYER_CONFIG_DIR/run/{runList}/{id}_*.json`
9. logging to `$DEPLOYER_LOGS_DIR/execute_{date}.log`
10. exclusive lock during run-list `execute` phase 

## Installation Instructions
`deployer` is designed to run as cronjob or as an on-demand cli directly on the host being managed. 
### 1. Build/Ship
The following command will install `deployer` on two hosts:
```shell
make build
targets=("3.88.103.159" "3.90.226.213")
for x in "${targets[@]}"
do
    scp ./build/deployer-linux root@${x}:deployer
done
```
### 2. Invoke Configuration (demo application)
```shell
## dry-run
./deployer execute -r hello-world

## real
./deployer execute -r hello-world --dryRun=false
```
### 3. Verify Demo Application Works 
```shell
curl -sv http://3.88.103.159
curl -sv http://3.90.226.213
```

## Examples
### Invoke Demo Configuration
```shell
root@ip-172-31-255-100:~# ./deployer execute -r hello-world --dryRun=false
2021/09/13 03:47:55 run_list.go:66: ... processing [0] Command
2021/09/13 03:47:55 directive_command.go:54: ... ( dryRun=false ) Command: apt-get update (depends on [])
2021/09/13 03:47:57 run_list.go:66: ... processing [1] OsPackage
2021/09/13 03:47:57 directive_os_package.go:90: ... ( dryRun=false ) OsPackage nginx action=install
2021/09/13 03:48:05 run_list.go:66: ... processing [2] OsPackage
2021/09/13 03:48:05 directive_os_package.go:90: ... ( dryRun=false ) OsPackage php-fpm action=install
2021/09/13 03:48:18 run_list.go:66: ... processing [3] GenerateFile
2021/09/13 03:48:18 directive_generate_file.go:121: ... ( dryRun=false ) GenerateFile: /etc/nginx/sites-available/hello-world.com action=create
2021/09/13 03:48:18 directive_generate_file.go:127: ... ( dryRun=false ) GenerateFile: /etc/nginx/sites-available/hello-world.com chmod=-rw-r-----
2021/09/13 03:48:18 directive_generate_file.go:133: ... ( dryRun=false ) GenerateFile: /etc/nginx/sites-available/hello-world.com chown=www-data:www-data (33:33)
2021/09/13 03:48:18 run_list.go:66: ... processing [4] Symlink
2021/09/13 03:48:18 directive_symlink.go:89: ... ( dryRun=false ) Symlink -> /etc/nginx/sites-available/hello-world.com: /etc/nginx/sites-enabled/hello-world.com action=create
2021/09/13 03:48:18 run_list.go:66: ... processing [5] Symlink
2021/09/13 03:48:18 directive_symlink.go:122: ... ( dryRun=false ) Symlink: /etc/nginx/sites-enabled/default action=delete
2021/09/13 03:48:18 run_list.go:66: ... processing [6] GenerateFile
2021/09/13 03:48:18 directive_generate_file.go:121: ... ( dryRun=false ) GenerateFile: /var/www/html/index.php action=create
2021/09/13 03:48:18 directive_generate_file.go:127: ... ( dryRun=false ) GenerateFile: /var/www/html/index.php chmod=-rw-r-----
2021/09/13 03:48:18 directive_generate_file.go:133: ... ( dryRun=false ) GenerateFile: /var/www/html/index.php chown=www-data:www-data (33:33)
2021/09/13 03:48:18 run_list.go:66: ... processing [7] Command
2021/09/13 03:48:18 directive_command.go:54: ... ( dryRun=false ) Command: nginx -t (depends on [1 2 3])
2021/09/13 03:48:19 run_list.go:66: ... processing [8] Command
2021/09/13 03:48:19 directive_command.go:54: ... ( dryRun=false ) Command: systemctl reload nginx (depends on [1 2 3])
root@ip-172-31-255-100:~# echo $?
0
```
### Idempotency
```shell
root@ip-172-31-255-100:~# ./deployer execute -r hello-world --dryRun=false
2021/09/13 03:49:22 run_list.go:66: ... processing [0] Command
2021/09/13 03:49:22 directive_command.go:54: ... ( dryRun=false ) Command: apt-get update (depends on [])
2021/09/13 03:49:24 run_list.go:66: ... processing [1] OsPackage
2021/09/13 03:49:24 directive_os_package.go:90: ... ( dryRun=false ) OsPackage nginx action=skip (already installed)
2021/09/13 03:49:24 run_list.go:66: ... processing [2] OsPackage
2021/09/13 03:49:24 directive_os_package.go:90: ... ( dryRun=false ) OsPackage php-fpm action=skip (already installed)
2021/09/13 03:49:24 run_list.go:66: ... processing [3] GenerateFile
2021/09/13 03:49:24 directive_generate_file.go:121: ... ( dryRun=false ) GenerateFile: /etc/nginx/sites-available/hello-world.com action=skip (file unchanged)
2021/09/13 03:49:24 directive_generate_file.go:127: ... ( dryRun=false ) GenerateFile: /etc/nginx/sites-available/hello-world.com chmod=-rw-r-----
2021/09/13 03:49:24 directive_generate_file.go:133: ... ( dryRun=false ) GenerateFile: /etc/nginx/sites-available/hello-world.com chown=www-data:www-data (33:33)
2021/09/13 03:49:24 run_list.go:66: ... processing [4] Symlink
2021/09/13 03:49:24 directive_symlink.go:89: ... ( dryRun=false ) Symlink -> /etc/nginx/sites-available/hello-world.com: /etc/nginx/sites-enabled/hello-world.com action=skip (already exists)
2021/09/13 03:49:24 run_list.go:66: ... processing [5] Symlink
2021/09/13 03:49:24 directive_symlink.go:122: ... ( dryRun=false ) Symlink: /etc/nginx/sites-enabled/default action=skip (already deleted)
2021/09/13 03:49:24 run_list.go:66: ... processing [6] GenerateFile
2021/09/13 03:49:24 directive_generate_file.go:121: ... ( dryRun=false ) GenerateFile: /var/www/html/index.php action=skip (file unchanged)
2021/09/13 03:49:24 directive_generate_file.go:127: ... ( dryRun=false ) GenerateFile: /var/www/html/index.php chmod=-rw-r-----
2021/09/13 03:49:24 directive_generate_file.go:133: ... ( dryRun=false ) GenerateFile: /var/www/html/index.php chown=www-data:www-data (33:33)
2021/09/13 03:49:24 run_list.go:66: ... processing [7] Command
2021/09/13 03:49:24 run_list.go:90: ... skipping Command(nginx -t) since none of its dependencies are modified
2021/09/13 03:49:24 run_list.go:66: ... processing [8] Command
2021/09/13 03:49:24 run_list.go:90: ... skipping Command(systemctl reload nginx) since none of its dependencies are modified
root@ip-172-31-255-100:~#
```
### Export Embedded Run-List
```shell
root@ip-172-31-255-100:~# ./deployer export -r hello-world
2021/09/13 03:49:59 sdk_export.go:25: ... exporting /etc/deployer/run/hello-world/0_Command.json
2021/09/13 03:49:59 sdk_export.go:25: ... exporting /etc/deployer/run/hello-world/1_OsPackage.json
2021/09/13 03:49:59 sdk_export.go:25: ... exporting /etc/deployer/run/hello-world/2_OsPackage.json
2021/09/13 03:49:59 sdk_export.go:25: ... exporting /etc/deployer/run/hello-world/3_GenerateFile.json
2021/09/13 03:49:59 sdk_export.go:25: ... exporting /etc/deployer/run/hello-world/4_Symlink.json
2021/09/13 03:49:59 sdk_export.go:25: ... exporting /etc/deployer/run/hello-world/5_Symlink.json
2021/09/13 03:49:59 sdk_export.go:25: ... exporting /etc/deployer/run/hello-world/6_GenerateFile.json
2021/09/13 03:49:59 sdk_export.go:25: ... exporting /etc/deployer/run/hello-world/7_Command.json
2021/09/13 03:49:59 sdk_export.go:25: ... exporting /etc/deployer/run/hello-world/8_Command.json
2021/09/13 03:49:59 sdk_export.go:49: ... exporting /etc/deployer/templates/hello-word.php.tmpl
2021/09/13 03:49:59 sdk_export.go:49: ... exporting /etc/deployer/templates/nginx-site.tmpl
```
once exported, you can invoke `deployer` and tell it to parse configuration directives from `/etc/deployer/*` (note `-s file` flag):
```shell
root@ip-172-31-255-100:~# ./deployer execute -r hello-world -s file --dryRun=false
2021/09/13 03:50:36 run_list.go:37: ... parsing 0_Command.json
2021/09/13 03:50:36 run_list.go:37: ... parsing 1_OsPackage.json
2021/09/13 03:50:36 run_list.go:37: ... parsing 2_OsPackage.json
2021/09/13 03:50:36 run_list.go:37: ... parsing 3_GenerateFile.json
2021/09/13 03:50:36 run_list.go:37: ... parsing 4_Symlink.json
2021/09/13 03:50:36 run_list.go:37: ... parsing 5_Symlink.json
2021/09/13 03:50:36 run_list.go:37: ... parsing 6_GenerateFile.json
2021/09/13 03:50:36 run_list.go:37: ... parsing 7_Command.json
2021/09/13 03:50:36 run_list.go:37: ... parsing 8_Command.json
2021/09/13 03:50:36 run_list.go:66: ... processing [0] Command
2021/09/13 03:50:36 directive_command.go:54: ... ( dryRun=false ) Command: apt-get update (depends on [])
2021/09/13 03:50:39 run_list.go:66: ... processing [1] OsPackage
2021/09/13 03:50:39 directive_os_package.go:90: ... ( dryRun=false ) OsPackage nginx action=skip (already installed)
2021/09/13 03:50:39 run_list.go:66: ... processing [2] OsPackage
2021/09/13 03:50:39 directive_os_package.go:90: ... ( dryRun=false ) OsPackage php-fpm action=skip (already installed)
2021/09/13 03:50:39 run_list.go:66: ... processing [3] GenerateFile
2021/09/13 03:50:39 directive_generate_file.go:121: ... ( dryRun=false ) GenerateFile: /etc/nginx/sites-available/hello-world.com action=skip (file unchanged)
2021/09/13 03:50:39 directive_generate_file.go:127: ... ( dryRun=false ) GenerateFile: /etc/nginx/sites-available/hello-world.com chmod=-rw-r-----
2021/09/13 03:50:39 directive_generate_file.go:133: ... ( dryRun=false ) GenerateFile: /etc/nginx/sites-available/hello-world.com chown=www-data:www-data (33:33)
2021/09/13 03:50:39 run_list.go:66: ... processing [4] Symlink
2021/09/13 03:50:39 directive_symlink.go:89: ... ( dryRun=false ) Symlink -> /etc/nginx/sites-available/hello-world.com: /etc/nginx/sites-enabled/hello-world.com action=skip (already exists)
2021/09/13 03:50:39 run_list.go:66: ... processing [5] Symlink
2021/09/13 03:50:39 directive_symlink.go:122: ... ( dryRun=false ) Symlink: /etc/nginx/sites-enabled/default action=skip (already deleted)
2021/09/13 03:50:39 run_list.go:66: ... processing [6] GenerateFile
2021/09/13 03:50:39 directive_generate_file.go:121: ... ( dryRun=false ) GenerateFile: /var/www/html/index.php action=skip (file unchanged)
2021/09/13 03:50:39 directive_generate_file.go:127: ... ( dryRun=false ) GenerateFile: /var/www/html/index.php chmod=-rw-r-----
2021/09/13 03:50:39 directive_generate_file.go:133: ... ( dryRun=false ) GenerateFile: /var/www/html/index.php chown=www-data:www-data (33:33)
2021/09/13 03:50:39 run_list.go:66: ... processing [7] Command
2021/09/13 03:50:39 run_list.go:90: ... skipping Command(nginx -t) since none of its dependencies are modified
2021/09/13 03:50:39 run_list.go:66: ... processing [8] Command
2021/09/13 03:50:39 run_list.go:90: ... skipping Command(systemctl reload nginx) since none of its dependencies are modified
root@ip-172-31-255-100:~#
```
### Sample Directives
OsPackage:
```shell
root@ip-172-31-255-100:~# cat /etc/deployer/run/hello-world/1_OsPackage.json
{
    "PkgName": "nginx",
    "Install": true,
    "Remove": false
}
```
GenerateFile:
```shell
root@ip-172-31-255-100:~# cat /etc/deployer/run/hello-world/3_GenerateFile.json
{
    "Path": "/etc/nginx/sites-available/hello-world.com",
    "Owner": "www-data",
    "Group": "www-data",
    "Mode": 416,
    "Template": "nginx-site",
    "Data": null,
    "Dependencies": null
}
```
File template for above ^^^:
```shell
root@ip-172-31-255-100:~# cat /etc/deployer/templates/nginx-site.tmpl
{{define "nginx-site"}}server {
        listen 80;
        root /var/www/html;
        index index.php index.html index.htm index.nginx-debian.html;
        server_name {{.GLOBAL_IP}};

        location / {
                try_files $uri $uri/ =404;
        }

        location ~ \.php$ {
                include snippets/fastcgi-php.conf;
                fastcgi_pass unix:/var/run/php/php7.2-fpm.sock;
        }

        location ~ /\.ht {
                deny all;
        }
}{{end}}
```

## Architecture
At the core, `deployer` defines and implements two interfaces:
1. `Directive`: a [directive](pkg/deployer/directive.go) is configuration manager for a specific resource (example: os-pkg, file, command, etc.).
2. `Runable`: a [run-list of directives](pkg/deployer/runable.go) listed in a [specific order](pkg/deployer/run_list_hello_world.go) of execution and linked with optional inter-dependency (for example, only restart service if config file is changed). 

There are two ways to define Directives and compose them into run-lists:
1. Via Deployer SDK functions: `RegisterDirective()` and `RegisterRunable()` [example see func NewSDK()](pkg/deployer/deployer.go)
2. Via file system configuration files stored in `$DEPLOYER_CONFIG_DIR/run/{ID}_{DIRECTIVE-NAME}.json` (see example below):
```shell
# ls -l /etc/deployer/run/hello-world/
total 36
-rw-r----- 1 root root  57 Sep 13 15:38 0_Command.json
-rw-r----- 1 root root  68 Sep 13 15:38 1_OsPackage.json
-rw-r----- 1 root root  70 Sep 13 15:38 2_OsPackage.json
-rw-r----- 1 root root 201 Sep 13 15:38 3_GenerateFile.json
-rw-r----- 1 root root 155 Sep 13 15:38 4_Symlink.json
-rw-r----- 1 root root 105 Sep 13 15:38 5_Symlink.json
-rw-r----- 1 root root 187 Sep 13 15:38 6_GenerateFile.json
-rw-r----- 1 root root  86 Sep 13 15:38 7_Command.json
-rw-r----- 1 root root 100 Sep 13 15:38 8_Command.json
```
### Demo
In addition to being able to register custom configuration management directives and compose them into run-list,
`deployer` implements an embedded demo run-list that installs `hello-world` application (`nginx`->`php-fpm`).
It's a great way to get familiar with its architecture:
#### Directives
* [GenerateFile](pkg/deployer/directive_generate_file.go)
* [OsPackage](pkg/deployer/directive_os_package.go)
* [Symlink](pkg/deployer/directive_symlink.go)
* [Command](pkg/deployer/directive_command.go)
#### Hello-World Run List
* [RunListHelloWorld](pkg/deployer/run_list_hello_world.go)
#### Execute Hello-World
```shell
## from embedded run-list
deployer execute -r hello-world --dryRun=false

## export to file system and then run off of that
deployer export -r hello-world
deployer execute -r hello-world -s file --dryRun=false
```
## 
END.
