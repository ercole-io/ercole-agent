// +build linux

// Copyright (c) 2020 Sorint.lab S.p.A.
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

package fetcher

import (
	"os/exec"
	"syscall"

	"github.com/ercole-io/ercole-agent/v2/logger"
)

// runCommandAs utility
func runCommandAs(log logger.Logger, u *User, commandName string, args ...string) (stdout, stderr []byte, exitCode int, err error) {
	cmd := exec.Command(commandName, args...)

	if u != nil {
		log.Debugf("runCommand [%v] with user [%v]", commandName, u)

		cmd.SysProcAttr = &syscall.SysProcAttr{}
		if err != nil {
			log.Errorf("Can't set process attributes at command [%v]", commandName)
			return nil, nil, -1, err
		}

		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: u.UID, Gid: u.GID}
	}

	stdout, err = cmd.Output()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			stderr = exitErr.Stderr
		} else {
			exitCode = -1
		}
	}

	return
}
