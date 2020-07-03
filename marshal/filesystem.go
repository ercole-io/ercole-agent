// Copyright (c) 2019 Sorint.lab S.p.A.
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

package marshal

import (
	"bufio"
	"strings"

	"github.com/ercole-io/ercole/model"
)

// Filesystems returns a list of Filesystem entries extracted
// from the filesystem fetcher command output.
// Filesystem output is a list of filesystem entries with positional attribute columns
// separated by one or more spaces
func Filesystems(cmdOutput []byte) []model.Filesystem {
	filesystems := []model.Filesystem{}

	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))

	for scanner.Scan() {
		line := scanner.Text()
		iter := newIter(strings.Fields(line))

		fs := model.Filesystem{}

		fs.Filesystem = strings.TrimSpace(iter())
		fs.Type = strings.TrimSpace(iter())
		fs.Size = TrimParseInt64(iter())
		fs.UsedSpace = TrimParseInt64(iter())
		fs.AvailableSpace = TrimParseInt64(iter())
		iter() // throw away used space percentage
		fs.MountedOn = strings.TrimSpace(iter())
		filesystems = append(filesystems, fs)
	}

	return filesystems
}
