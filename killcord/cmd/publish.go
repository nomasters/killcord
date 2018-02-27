package cmd

import (
	"fmt"
	"os"

	"github.com/nomasters/killcord"
	"github.com/spf13/cobra"
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publishes the decryption key to the contract",
	Long:  `Publishes the decryption key to the contract`,
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

		session := killcord.New()
		session.Config = config
		session.Options = setOptionsFromEnv()
		session.Init()

		if err := session.PublishKey(); err != nil {
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
	rootCmd.AddCommand(publishCmd)
}
