package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/nomasters/killcord"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cfgFileExists bool
var devMode bool

var rootCmd = &cobra.Command{
	Use:   "killcord",
	Short: "A censorship resistant dead man's switch",
	Long: `killcord is designed for a project owner to be able to automatically release 
a decryption key in the circumstance that the owner stops checking-in after a predefined 
window of time. This allows an owner to release a secret data dump to the public in the 
event of death or disappearance.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./killcord.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// get the configuration file from the current directory
		viper.AddConfigPath(".")
		viper.SetConfigName("killcord")
	}

	// support ENV for killcord settings
	viper.SetEnvPrefix("killcord")
	viper.BindEnv("payload.secret", "KILLCORD_PAYLOAD_SECRET")
	viper.BindEnv("contract.id", "KILLCORD_CONTRACT_ID")
	viper.BindEnv("contract.publisher.address", "KILLCORD_CONTRACT_PUBLISHER_ADDRESS")
	viper.BindEnv("contract.publisher.password", "KILLCORD_CONTRACT_PUBLISHER_PASSWORD")
	viper.BindEnv("contract.publisher.keystore", "KILLCORD_CONTRACT_PUBLISHER_KEYSTORE")
	viper.BindEnv("contract.rpcUrl", "KILLCORD_CONTRACT_RPCURL")
	viper.BindEnv("payload.rpcUrl", "KILLCORD_PAYLOAD_RPCURL")
	viper.BindEnv("publisher.warningThreshold", "KILLCORD_PUBLISHER_WARNING_THRESHOLD")
	viper.BindEnv("publisher.publishThreshold", "KILLCORD_PUBLISHER_PUBLISH_THRESHOLD")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		cfgFileExists = true
	}
}

func setOptionsFromEnv() killcord.ProjectOptions {
	opts := killcord.ProjectOptions{}
	opts.DevMode = devMode
	if x, ok := viper.Get("payload.id").(string); ok {
		opts.Payload.ID = x
	}
	if x, ok := viper.Get("payload.secret").(string); ok {
		opts.Payload.Secret = x
	}
	if x, ok := viper.Get("publisher.warningThreshold").(int64); ok {
		opts.Publisher.WarningThreshold = x
	}
	if x, ok := viper.Get("publisher.publishThreshold").(int64); ok {
		opts.Publisher.PublishThreshold = x
	}
	if x, ok := viper.Get("contract.id").(string); ok {
		opts.Contract.ID = x
	}
	if x, ok := viper.Get("contract.publisher.address").(string); ok {
		opts.Contract.Publisher.Address = x
	}
	if x, ok := viper.Get("contract.publisher.password").(string); ok {
		opts.Contract.Publisher.Password = x
	}
	if x, ok := viper.Get("contract.publisher.keystore").(string); ok {
		opts.Contract.Publisher.KeyStore = x
	}
	if x, ok := viper.Get("contract.rpcUrl").(string); ok {
		opts.Contract.RPCURL = x
	}
	if x, ok := viper.Get("payload.rpcUrl").(string); ok {
		opts.Payload.RPCURL = x
	}
	return opts
}

func getConfigFile() (killcord.ProjectConfig, error) {
	var config killcord.ProjectConfig
	if _, err := toml.DecodeFile("killcord.toml", &config); err != nil {
		return config, err
	}
	return config, nil
}

func updateConfigFile(config killcord.ProjectConfig) error {
	fo, err := os.OpenFile("killcord.toml", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	w := bufio.NewWriter(fo)
	c := toml.NewEncoder(w)
	if err := c.Encode(&config); err != nil {
		return err
	}
	if err := viper.ReadInConfig(); err == nil {
		cfgFileExists = true
	} else {
		return err
	}
	return nil
}
