package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/gofrs/flock"
	"github.com/vmogilev/deployer/internal/env"
)

type job struct {
	verbose     bool
	forceUnLock bool
	log         *log.Logger
	vars        *env.Vars
	lockHandle  *flock.Flock
}

func (j *job) lock() {
	j.lockHandle = flock.New(j.vars.DeployerLockFile)
	if j.forceUnLock {
		j.unlock()
	}

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
