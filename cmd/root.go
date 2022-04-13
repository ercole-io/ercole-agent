// Copyright (c) 2022 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
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
	Run: func(cmd *cobra.Command, args []string) {
		serve(new(program))
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
