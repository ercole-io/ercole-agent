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
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var olaCmd = &cobra.Command{
	Use:   "ola",
	Short: "Oracle License Assessment",
	Long:  `Manage Oracle License Assessment`,
	Run: func(cmd *cobra.Command, args []string) {
		hostname, errHost := os.Hostname()
		if errHost != nil {
			fmt.Fprintf(os.Stderr, "%s\n", errHost)
			log.Printf("Error: %s\n", errHost)
			os.Exit(1)
		}

		fileLogPath := fmt.Sprintf("%s/ercole-agent-ola-%s-%s.log", os.TempDir(), hostname, time.Now().Local().Format("06-01-02_15-04-05"))

		f, err := os.Create(fileLogPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			log.Printf("Error: %s\n", errHost)
			os.Exit(1)
		}
		defer f.Close()
		log.SetOutput(f)

		fmt.Printf("Start ola %s \n", time.Now().Local().Format("06-01-02_15-04-05"))
		log.Printf("Start ola %s \n", time.Now().Local().Format("06-01-02_15-04-05"))

		if err := checkOlaOratab(configuration.Features.OracleDatabase.Oratab); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			log.Printf("Error: %s\n", errHost)
			os.Exit(1)
		}

		if err := suggestOlaEntry(configuration.Features.OracleDatabase.Oratab); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			log.Printf("Error: %s\n", errHost)
			os.Exit(1)
		}

		if err := manageOlaVerbose(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			log.Printf("Error: %s\n", errHost)
			os.Exit(1)
		}

		fmt.Printf("End ola %s \n", time.Now().Local().Format("06-01-02_15-04-05"))
		log.Printf("End ola %s \n", time.Now().Local().Format("06-01-02_15-04-05"))
	},
}

func init() {
	rootCmd.AddCommand(olaCmd)
}
func checkOlaOratab(filePath string) error {
	fmt.Printf("Check oratab step: %s \n", filePath)
	log.Printf("Check oratab step: %s \n", filePath)

	if _, err := os.Stat(filePath); err != nil {
		fmt.Printf("Check oratab error: %s \n", err)
		log.Printf("Check oratab error: %s \n", err)

		return err
	}

	return nil
}

func suggestOlaEntry(filePath string) error {
	pathScript := "fetch/linux/suggest_oratab.sh"

	cmd := exec.CommandContext(context.TODO(), pathScript, filePath)
	fmt.Printf("Suggest entry step: %s %s \n", pathScript, filePath)
	log.Printf("Suggest entry step: %s %s \n", pathScript, filePath)

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Printf("Suggest entry error: %s \n", err)
		log.Printf("Suggest entry error: %s \n", err)

		return err
	}

	stringstdout := string(stdout)

	if stringstdout != "" {
		fmt.Printf("Suggested entry result:\n %s \n", stringstdout)
		log.Printf("Suggested entry result:\n %s \n", stringstdout)

		fmt.Println("Do you want to add the suggested entries? (Y/N)")
		log.Println("Asked: `Do you want to add the suggested entries? (Y/N)`")

		reader := bufio.NewReader(os.Stdin)
		char, _, errReader := reader.ReadRune()

		if errReader != nil {
			fmt.Printf("Reader error: %s \n", errReader)
			log.Printf("Reader error: %s \n", errReader)

			return errReader
		}

		resp := strings.ToLower(string(char))

		switch resp {
		case "y":
			err := changeOlaOratab(configuration.Features.OracleDatabase.Oratab, stringstdout)
			if err != nil {
				fmt.Printf("Change oratab file error: %s \n", err)
				log.Printf("Change oratab file error: %s \n", err)

				return err
			}

			fmt.Println("Changed oratab entries")
			log.Println("Changed oratab entries")

		case "n":
			fmt.Println("NOT changed oratab entries")
			log.Println("NOT changed oratab entries")
		default:
			log.Println("Value not valid: only values Y and N are allowed")

			return fmt.Errorf("Value not valid: only values Y and N are allowed")
		}
	} else {
		stringstdout = "No suggestions: it's all ok"

		fmt.Printf("Suggested entry result: %s \n", stringstdout)
		log.Printf("Suggested entry result: %s \n", stringstdout)
	}

	return nil
}

func manageOlaVerbose() error {
	pathScript := "fetch/linux/exec_verbose.sh"

	cmd := exec.CommandContext(context.TODO(), pathScript)
	fmt.Printf("Manage verbose step: %s \n", pathScript)
	log.Printf("Manage verbose step: %s \n", pathScript)

	pipe, errOut := cmd.StdoutPipe()
	if errOut != nil {
		fmt.Printf("Manage verbose error: %s \n", errOut)
		log.Printf("Manage verbose error: %s \n", errOut)

		return errOut
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Manage verbose error: %s \n", err)
		log.Printf("Manage verbose error: %s \n", err)

		return err
	}

	reader := bufio.NewReader(pipe)
	line, err := reader.ReadString('\n')

	for err == nil {
		fmt.Print(line)
		line, err = reader.ReadString('\n')
	}

	return nil
}

func changeOlaOratab(filePath string, text string) error {
	fmt.Println("Manage change oratab file step")
	log.Println("Manage change oratab file step")

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Printf("Open oratab file error: %s \n", err)
		log.Printf("Open oratab file error: %s \n", err)

		return err
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		fmt.Printf("Write oratab file error: %s \n", err)
		log.Printf("Write oratab file error: %s \n", err)

		return err
	}

	fmt.Printf("Hostdata pretty-printed on file: %s \n", filePath)
	log.Printf("Hostdata pretty-printed on file: %s \n", filePath)

	return nil
}
