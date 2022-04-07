package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/ercole-io/ercole-agent/v2/config"
	"github.com/ercole-io/ercole-agent/v2/logger"
	"github.com/spf13/cobra"
)

var configuration config.Configuration
var verbose bool
var extraConfigFile string

var rootCmd = &cobra.Command{
	Use:   "ercole-agent",
	Short: "ercole-agent",
	Long:  `ercole-agent`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		confLog, err := logger.NewLogger("CONFIG", logger.LogVerbosely(verbose))
		if err != nil {
			log.Fatal("Can't initialize CONFIG logger: ", err)
		}

		extraConfigFile = strings.TrimSpace(extraConfigFile)

		if len(extraConfigFile) > 0 && !fileExists(extraConfigFile) {
			log.Fatalf("Configuration file not found: %s", extraConfigFile)
		}

		configuration = config.ReadConfig(confLog, extraConfigFile)

		if verbose {
			configuration.Verbose = verbose
		}
	},
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&extraConfigFile, "config", "c", "", "Configuration file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
}
