package cmd

import (
	"fmt"
	"os"

	"github.com/nomasters/killcord"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a killcord project",
	Long: `
The killcord init command creates the project folder structure required
to own a killcord project. If you plan on releasing a payload to the 
public in the event of death or disappearance, start here.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Ensure no configuration file exists.
		// If one does, exit with message.
		if cfgFileExists == true {
			fmt.Println("killcord.toml discovered, killcord is already initialized")
			os.Exit(1)
		}

		// Set default options for a public killcord. This is a structural placeholder
		// for future features in which Payload and Contract Providers can be switched
		// out independently.
		opts := setOptionsFromEnv()
		opts.Type = "owner"
		opts.Audience = "public"
		opts.Payload.Provider = "ipfs"
		opts.Contract.Provider = "ethereum"

		session := killcord.New()
		session.Options = opts
		session.Init()

		// Start a new project with the options provided and return ProjectConfig
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
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&devMode, "dev", false, "initialize project to run in dev mode")
}
