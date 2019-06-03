package cmd

import (
	"strings"

	"github.com/open-identity/utils/logrusx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version = "master"
	Date    = "undefined"
	Commit  = "undefined"
)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "roles-manager",
}

var logger = new(logrus.Logger)

// Execute adds all child commands to the root command sets flags appropriately.
// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	*logger = *logrusx.New()
}
