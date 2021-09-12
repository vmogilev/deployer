package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vmogilev/deployer/internal/env"
	"github.com/vmogilev/deployer/pkg/deployer"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "exports embedded run-list to ${DEPLOYER_CONFIG_DIR}/{runList}",
	Long: `
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

		sdk := deployer.NewSDK(j.log, &deployer.SDKInput{
			Verbose: j.verbose,
			Vars:    vars,
		})
		if err := sdk.Export(viper.GetString("runList")); err != nil {
			j.abort(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().BoolP("verbose", "v", false, "verbose")
	exportCmd.Flags().StringP("runList", "r", "", "run-list name to export")
}
