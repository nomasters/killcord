package cmd

import (
	"fmt"

	"github.com/nomasters/killcord"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the killcord version",
	Long:  `Prints the killcord version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(killcord.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
