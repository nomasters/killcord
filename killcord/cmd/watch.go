package cmd

import (
	"fmt"
	"os"

	"github.com/nomasters/killcord"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Initializes a watcher project",
	Long: `The watch command is used to initialize a project to watch an existing project.
The watcher downloads the payload and can run the status and decrypt commands on a payload.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Ensure no configuration file exists.
		// If one does, exit with message.
		if cfgFileExists == true {
			fmt.Println("killcord.toml discovered, killcord is already initialized")
			os.Exit(1)
		}

		// Set default options for a killcord watcher project. This is a structural placeholder
		// for future features in which Payload and Contract Providers can be switched
		// out independently.

		session := killcord.New()
		opts := setOptionsFromEnv()
		opts.Type = "watcher"
		opts.Payload.Provider = "ipfs"
		opts.Contract.Provider = "ethereum"
		opts.Contract.ID = args[0]
		session.Options = opts
		session.Init()
		if err := session.NewProject(); err != nil {
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
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().BoolVar(&devMode, "dev", false, "run in dev mode")
}
