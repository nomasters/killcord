package cmd

import (
	"fmt"
	"os"

	"github.com/nomasters/killcord"
	"github.com/spf13/cobra"
)

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Manage a contract",
	Long: `More options will become available, but currently this comamand only 
accepts the --kill flag which is used to kill an existing contract`,
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
		kill, err := cmd.Flags().GetBool("kill")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		session := killcord.New()
		session.Config = config
		session.Options = setOptionsFromEnv()
		session.Init()

		if kill == true {
			if err := session.KillContract(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		// write config to disk
		if err := updateConfigFile(session.Config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var deployContractCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a contract",
	Long: `This command deploys a smart contract and writes all the configuration 
details to the config file. Deploying a contract requires that the owner account 
has enough currency provided to fund the the deployment of a contract.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Ensure no configuration file exists.
		// If one does, exit with message.
		if cfgFileExists == false {
			fmt.Println("no killcord.toml found, run `killcord init` to start a new project")
			os.Exit(1)
		}
		// Set default status options for a public killcord.
		opts := setOptionsFromEnv()
		config, err := getConfigFile()
		if err != nil {
			fmt.Println("something went wrong with reading the config file, exiting.")
			os.Exit(1)
		}

		session := killcord.New()
		session.Options = opts
		session.Config = config
		session.Init()
		if err := session.DeployContract(); err != nil {
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
	contractCmd.Flags().BoolP("kill", "k", false, "kill contract")
	rootCmd.AddCommand(contractCmd)
	contractCmd.AddCommand(deployContractCmd)
}
