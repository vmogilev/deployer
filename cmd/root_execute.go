package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vmogilev/deployer/internal/env"
	"github.com/vmogilev/deployer/pkg/deployer"
)

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "executes deployer daemon and applies --runList directives from --source (see below)",
	Long: `

Executes deployer daemon and applies --runList directives from --source:

   --source=build :: (default) directives and templates are sourced from the binary (embedded) 
   --source=file  :: directives and templates are sourced from ${DEPLOYER_CONFIG_DIR}/run/{runList}

By default, it runs in dry-run mode and will only print the actions it's about to take.  
To apply the changes, add --dryRun=false flag.

EXAMPLES:

## list available run-lists from build (embedded)
deployer execute -l

## list available run-lists from from ${DEPLOYER_CONFIG_DIR}/run/
deployer execute -l -s file

## dry run all directives for hello-world runList from build (embedded)
deployer execute -r hello-world

## apply all directives for hello-world runList from build (embedded)
deployer execute -r hello-world --dryRun=false

## apply all directives for hello-world runList from ${DEPLOYER_CONFIG_DIR}/run/hello-world
deployer execute -r hello-world -s file --dryRun=false


ENV VARs:
	DEPLOYER_CONFIG_DIR	  default:/etc/deployer
	DEPLOYER_LOGS_DIR     default:/var/log/deployer
	DEPLOYER_LOCK_FILE    default:/var/lock/deployer.lock
	LOCK_TIMEOUT_SECONDS  default:15


Concurrency:
	Deployer can safely run as daemon/cronjob since it takes an exclusive lock via ${DEPLOYER_LOCK_FILE}
`,
	Run: func(cmd *cobra.Command, args []string) {
		log := log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags|log.LUTC)

		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			log.Fatalln(err)
		}

		vars, err := env.LoadVars()
		if err != nil {
			log.Fatalf("error parsing ENV VARs: %v", err)
		}

		j := &job{
			vars:    vars,
			verbose: viper.GetBool("verbose"),
			log:     log,
		}
		j.init()

		// setup multi-writer to log to stdout and a logfile in DEPLOYER_LOGS_DIR
		ln := fmt.Sprintf("execute_%s.log", time.Now().UTC().Format("2006-01-02-150405"))
		fp := filepath.Join(vars.DeployerLogsDir, ln)
		logFile, err := os.OpenFile(fp, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Fatalf("can't create logfile %s: %v", fp, err)
		}
		// this ignores the Close() error
		defer logFile.Close()
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)

		j.lock()
		defer j.unlock()

		sdk, err := deployer.NewSDK(j.log, &deployer.SDKInput{
			Verbose: j.verbose,
			Vars:    vars,
		})
		if err != nil {
			j.abort(err.Error())
		}

		var exe func(ctx context.Context, in *deployer.ExecuteInput) error
		switch viper.GetString("source") {
		case "build":
			exe = sdk.ExecuteFromBuild
		case "file":
			exe = sdk.ExecuteFromFS
		}

		in := &deployer.ExecuteInput{
			RunList:  viper.GetString("runList"),
			DryRun:   viper.GetBool("dryRun"),
			Force:    viper.GetBool("forceReconcile"),
			ListOnly: viper.GetBool("list"),
		}
		if err := in.Validate(); err != nil {
			j.abort(err.Error())
		}
		if err := exe(context.Background(), in); err != nil {
			j.abort(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
	executeCmd.Flags().BoolP("verbose", "v", false, "verbose")
	executeCmd.Flags().Bool("dryRun", true, "dryRun")
	executeCmd.Flags().Bool("forceReconcile", false, "force reconcile")
	executeCmd.Flags().StringP("source", "s", "build", "build|file (build run-lists are embedded into binary, file run-lists are pulled from $DEPLOYER_CONFIG_DIR)")
	executeCmd.Flags().StringP("runList", "r", "", "run-list name to execute (see --source for where they are defined)")
	executeCmd.Flags().BoolP("list", "l", false, "list available run-lists from --source")
}
