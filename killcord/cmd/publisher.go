package cmd

import (
	"fmt"
	"os"

	"github.com/nomasters/killcord"
	"github.com/spf13/cobra"
)

var publisherCmd = &cobra.Command{
	Use:   "publisher",
	Short: "Manage the Publisher",
	Long: `This Command is a placeholder for future publisher commands.
Currently the only publisher command is "killcord publisher run".`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var runPublisherCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the Publisher threshold check.",
	Long: `Publisher Run is the primary interface for automating the publisher
threshold check. This command is configured to use either a config file or
ENV settings to read a smart contract's lastCheckin value and evaluate it against
the publisher threshold. If the threshold exceeded, the command automatically 
publishes the secret to the smart contract. If the secret is already published,
this step is skipped.`,
	Run: func(cmd *cobra.Command, args []string) {

		// get config and opts
		config, _ := getConfigFile() // if no config file is Found this is fine
		opts := setOptionsFromEnv()

		session := killcord.New()
		session.Options = opts
		session.Config = config
		session.Init()

		if err := session.RunPublisher(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(publisherCmd)
	publisherCmd.AddCommand(runPublisherCmd)
	runPublisherCmd.Flags().BoolVar(&devMode, "dev", false, "run in dev mode")
}
