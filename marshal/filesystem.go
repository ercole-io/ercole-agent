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
	"bytes"
	"strings"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/hashicorp/go-multierror"
)

// Filesystems returns a list of Filesystem entries extracted
// from the filesystem fetcher command output.
// Filesystem output is a list of filesystem entries with positional attribute columns
// separated by one or more spaces
func Filesystems(cmdOutput []byte) ([]model.Filesystem, error) {
	filesystems := []model.Filesystem{}

	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))
	var merr error
	var err error

	for scanner.Scan() {
		line := scanner.Text()
		iter := NewIter(strings.Fields(line))

		fs := model.Filesystem{}

		fs.Filesystem = strings.TrimSpace(iter())
		fs.Type = strings.TrimSpace(iter())
		if fs.Size, err = TrimParseInt64(iter()); err != nil {
			merr = multierror.Append(merr, err)
		}

		if fs.UsedSpace, err = TrimParseInt64(iter()); err != nil {
			merr = multierror.Append(merr, err)
		}

		if fs.AvailableSpace, err = TrimParseInt64(iter()); err != nil {
			merr = multierror.Append(merr, err)
		}

		iter() // throw away used space percentage
		fs.MountedOn = strings.TrimSpace(iter())
		filesystems = append(filesystems, fs)
	}

	if merr != nil {
		return nil, merr
	}
	return filesystems, nil
}
