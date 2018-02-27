package cmd

import (
	"fmt"
	"os"

	"github.com/nomasters/killcord"
	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypts a payload",
	Long: `The decrypt command uses the settings stored in the configuration file to decrypt
a payload stored in the payload/encrypted directory.`,
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
		session.Init()

		if err := session.Decrypt(); err != nil {
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
	rootCmd.AddCommand(decryptCmd)
}
