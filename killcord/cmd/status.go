package cmd

import (
	"fmt"
	"os"

	"github.com/nomasters/killcord"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the project",
	Long: `
Status is used to check important status information about
an existing project.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Ensure no configuration file exists.
		// If one does, exit with message.
		if cfgFileExists == false {
			fmt.Println("no killcord.toml found, run `killcord init` to start a new project")
			os.Exit(1)
		}

		config, err := getConfigFile()
		if err != nil {
			fmt.Println("something went wrong with reading the config file, exiting.")
			os.Exit(1)
		}
		opts := setOptionsFromEnv()

		session := killcord.New()
		session.Config = config
		session.Options = opts
		session.Init()

		if err := session.GetStatus(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// write config to disk
		if err := updateConfigFile(session.Config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
