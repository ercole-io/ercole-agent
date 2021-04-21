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

// +build !windows

package config

import (
	"fmt"
	"os"
	"syscall"
)

func isDirectoryWritable(path string) (isWritable bool, err error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("Path doesn't exist: %w", err)
	}

	if !info.IsDir() {
		return false, fmt.Errorf("Path isn't a directory")
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return false, fmt.Errorf("Write permission bit is not set on this file for user")
	}

	var stat syscall.Stat_t
	if err = syscall.Stat(path, &stat); err != nil {
		return false, fmt.Errorf("Unable to get stat: %w", err)
	}

	if uint32(os.Geteuid()) != stat.Uid {
		return false, fmt.Errorf("User doesn't have permission to write to this directory")
	}

	return true, nil
}
