package cmd

import (
	"fmt"
	"os"

	"github.com/nomasters/killcord"
	"github.com/spf13/cobra"
)

var checkinCmd = &cobra.Command{
	Use:   "checkin",
	Short: "Checkin with current timestamp",
	Long: `The checkin command is used by the owner account to update the
checking timestamp for a contract.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		if err := session.CheckIn(); err != nil {
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
	rootCmd.AddCommand(checkinCmd)
}
