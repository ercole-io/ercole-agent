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
	"fmt"
	"os"

	"github.com/ercole-io/ercole-agent/v2/builder"
	"github.com/ercole-io/ercole-agent/v2/logger"
	"github.com/spf13/cobra"
)

var verboseCmd = &cobra.Command{
	Use:   "verbose",
	Short: "Write json file with hostdata",
	Long:  `Write json file with hostdata`,
	Run: func(cmd *cobra.Command, args []string) {

		if err := manageHostDataVerbose(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(verboseCmd)
}

func manageHostDataVerbose() error {
	opts := make([]logger.LoggerOption, 0)
	opts = append(opts, logger.LogVerbosely(true))

	logBuilData, err := logger.NewLogger("AGENT OLA", opts...)
	if err != nil {
		logBuilData.Fatal("Can't initialize AGENT logger: ", err)

		return err
	}

	hostData := builder.BuildData(configuration, logBuilData)

	hostData.AgentVersion = version
	hostData.SchemaVersion = hostDataSchemaVersion
	hostData.Period = configuration.Period
	hostData.Tags = []string{}

	writeHostDataOnTmpFile(hostData, logBuilData)

	return nil
}
