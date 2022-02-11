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

	"golang.org/x/sys/unix"
)

func checkDirectoryIsWritable(path string) (err error) {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("Path doesn't exist: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("Path isn't a directory")
	}

	if err := unix.Access(path, unix.W_OK); err != nil {
		return fmt.Errorf("User has no write permission: %w", err)
	}

	if err := unix.Access(path, unix.X_OK); err != nil {
		return fmt.Errorf("User has no execute permission: %w", err)
	}

	return nil
}
