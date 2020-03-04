package builder

import (
	"bytes"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ercole-io/ercole-agent/config"
)

func fetcher(configuration config.Configuration, fetcherName string, params ...string) []byte {
	var (
		cmd    *exec.Cmd
		err    error
		psexe  string
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	baseDir := config.GetBaseDir()

	if runtime.GOOS == "windows" {
		psexe, err = exec.LookPath("powershell.exe")
		if err != nil {
			log.Fatal(psexe)
		}
		if configuration.ForcePwshVersion == "0" {
			params = append([]string{"-ExecutionPolicy", "Bypass", "-File", baseDir + "\\fetch\\win.ps1", "-s", fetcherName}, params...)
		} else {
			params = append([]string{"-version", configuration.ForcePwshVersion, "-ExecutionPolicy", "Bypass", "-File", baseDir + "\\fetch\\win.ps1", "-s", fetcherName}, params...)
		}
		log.Println("Fetching " + psexe + " " + strings.Join(params, " "))

		cmd = exec.Command(psexe, params...)
	} else {
		log.Println("Fetching " + baseDir + "/fetch/" + fetcherName + " " + strings.Join(params, " "))
		cmd = exec.Command(baseDir+"/fetch/"+fetcherName, params...)
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if len(stderr.Bytes()) > 0 {
		log.Print(string(stderr.Bytes()))
	}

	if err != nil {
		if fetcherName != "dbstatus" {
			log.Fatal(err)
		} else {
			return []byte("UNREACHABLE") // fallback
		}
	}

	return stdout.Bytes()
}
