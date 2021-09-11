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

const provisionerLock = "provisioner"

// executeCmd represents the start command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "executes deployer daemon and applies all pending run-lists in ${DEPLOYER_CONFIG_DIR}/run-list",
	Long: `

Deployer encodes run-list directives in its codebase and applies them by iterating over a run-list.
By default, it runs in dry-run mode and will only print the actions it's about to take.  
To apply the changes, add --dryRun=false flag.

EXAMPLE:

## dry run all pending changes
deployer execute

## apply all pending changes
deployer execute --dryRun=false
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

		j := &job{
			vars:    vars,
			verbose: viper.GetBool("verbose"),
			log:     log,
		}
		j.lock()
		defer j.unlock()
		j.init()

		sdk := deployer.NewSDK(j.log, &deployer.SDKInput{
			Verbose: j.verbose,
			Vars:    vars,
		})
		if err := sdk.Execute(context.Background(), &deployer.ExecuteInput{
			DryRun: viper.GetBool("dryRun"),
			Force:  viper.GetBool("forceReconcile"),
		}); err != nil {
			j.abort(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
	executeCmd.Flags().BoolP("verbose", "v", false, "verbose")
	executeCmd.Flags().Bool("dryRun", true, "dryRun")
	executeCmd.Flags().Bool("forceReconcile", false, "force reconcile")
}
