package cmd

import (
	"fmt"
	"os"

	"github.com/marcusyip/golang-wire-mongo/cmd/migrate"
	"github.com/marcusyip/golang-wire-mongo/config"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			conf := config.ProvideConfig()
			app, err := BuildApp(conf)
			if err != nil {
				panic(err)
			}
			app.Start()
		},
	}
	cmd.AddCommand(
		migrate.NewMigrateCommand(),
	)
	return cmd
}

func onInitialize() {
	// loadConfig()
	// repos.Init()

	// if config.Get().App.SeedSampleData {
	// 	sampledata.Seed()
	// }

	// httputil.ErrPageUrl = config.Get().WorkApps.AccountsURL + "/error"
}

func init() {
	cobra.OnInitialize(onInitialize)
	rootCmd = newRootCommand()
}
