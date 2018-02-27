package cmd

import (
	"fmt"
	"os"

	"github.com/nomasters/killcord"
	"github.com/spf13/cobra"
)

var payloadCmd = &cobra.Command{
	Use:   "payload",
	Short: "Manage the payload",
	Long: `This is a placeholder for future work. To use the payload command
you will need to run "killcord payload deploy"`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var payloadDeploy = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a payload to a storage endpoint",
	Long: `Payload deploy deploys a payload to a storage endpoint and
registers that enpoint to the smart contract.`,
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

		if err := session.DeployPayload(); err != nil {
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
	payloadCmd.AddCommand(payloadDeploy)
	rootCmd.AddCommand(payloadCmd)
}
