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
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var oratabCmd = &cobra.Command{
	Use:   "check-oratab",
	Short: "Manage oratab file",
	Long:  `Set oratab file with expected value`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkOratab(configuration.Features.OracleDatabase.Oratab); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		if err := suggestEntry(configuration.Features.OracleDatabase.Oratab); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(oratabCmd)
}
func checkOratab(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		return err
	}

	return nil
}

func suggestEntry(filePath string) error {
	cmd := exec.CommandContext(context.TODO(), "fetch/linux/suggest_oratab.sh", filePath)

	stdout, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Print(string(stdout))

	return nil
}
