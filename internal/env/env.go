package env

import (
	"github.com/kelseyhightower/envconfig"
)

// Vars - list of environmental variables we expect
type Vars struct {
	DeployerConfigDir  string `envconfig:"DEPLOYER_CONFIG_DIR" default:"/etc/deployer"`
	DeployerLogsDir    string `envconfig:"DEPLOYER_LOGS_DIR" default:"/var/log/deployer"`
	DeployerLockFile   string `envconfig:"DEPLOYER_LOCK_FILE" default:"/var/lock/deployer.lock"`
	LockTimeoutSeconds int    `envconfig:"LOCK_TIMEOUT_SECONDS" default:"15"`
}

// LoadVars - loads env vars
func LoadVars() (*Vars, error) {
	var ev Vars
	noPrefix := ""
	if err := envconfig.Process(noPrefix, &ev); err != nil {
		return nil, err
	}
	return &ev, nil
}
