package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/gofrs/flock"
	"github.com/vmogilev/deployer/internal/env"
	"github.com/vmogilev/deployer/pkg/deployer"
)

type job struct {
	verbose    bool
	log        *log.Logger
	vars       *env.Vars
	lockHandle *flock.Flock
}

func (j *job) mkdirAll(path string) {
	if err := os.MkdirAll(path, 0700); err != nil {
		j.abort(err.Error())
	}
}

func (j *job) init() {
	j.mkdirAll(j.vars.DeployerLogsDir)
	j.mkdirAll(filepath.Join(j.vars.DeployerConfigDir, deployer.DirNameRunList))
	j.mkdirAll(filepath.Join(j.vars.DeployerConfigDir, deployer.DirNameTemplates))
	j.mkdirAll(filepath.Join(j.vars.DeployerConfigDir, deployer.DirNameCache))
	j.mkdirAll(filepath.Dir(j.vars.DeployerLockFile))
}

func (j *job) lock() {
	j.lockHandle = flock.New(j.vars.DeployerLockFile)

	var cnt int
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		if cnt >= j.vars.LockTimeoutSeconds {
			ticker.Stop()
			j.log.Fatalf("lock timeout after %d seconds", j.vars.LockTimeoutSeconds)
		}
		ok, err := j.lockHandle.TryLock()
		if err != nil {
			ticker.Stop()
			j.log.Fatalf("can't create lock handle: %v", err)
		}
		if ok {
			ticker.Stop()
			return
		}
		j.log.Println("spin getting a lock ...")
		cnt++
	}
}

func (j *job) unlock() {
	if j.lockHandle == nil {
		return
	}
	if err := j.lockHandle.Unlock(); err != nil {
		j.log.Println(err)
	}
}

func (j *job) whoAmI() string {
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s@%s", user, host)
}

func (j *job) abort(mess string) {
	j.unlock()
	j.log.Fatalln(mess)
}
