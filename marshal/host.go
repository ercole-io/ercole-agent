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
	"strings"

	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
	"github.com/hashicorp/go-multierror"
)

// Host returns a Host struct from the output of the host
// fetcher command. Host fields output is in key: value format separated by a newline
func Host(cmdOutput []byte) (*model.Host, error) {
	data := parseKeyValueColonSeparated(cmdOutput)

	var m model.Host

	var merr, err error

	m.Hostname = strings.TrimSpace(data["Hostname"])
	m.CPUModel = strings.TrimSpace(data["CPUModel"])
	m.CPUFrequency = strings.TrimSpace(data["CPUFrequency"])

	if m.CPUSockets, err = TrimParseInt(data["CPUSockets"]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if m.CPUCores, err = TrimParseInt(data["CPUCores"]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if m.CPUThreads, err = TrimParseInt(data["CPUThreads"]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if m.ThreadsPerCore, err = TrimParseInt(data["ThreadsPerCore"]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if m.CoresPerSocket, err = TrimRoundFloatToInt(data["CoresPerSocket"]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	m.HardwareAbstraction = strings.TrimSpace(data["HardwareAbstraction"])
	m.HardwareAbstractionTechnology = strings.TrimSpace(data["HardwareAbstractionTechnology"])
	m.Kernel = strings.TrimSpace(data["Kernel"])
	m.KernelVersion = strings.TrimSpace(data["KernelVersion"])
	m.OS = strings.TrimSpace(data["OS"])
	m.OSVersion = strings.TrimSpace(data["OSVersion"])

	if m.MemoryTotal, err = TrimParseFloat64(data["MemoryTotal"]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if m.SwapTotal, err = TrimParseFloat64(data["SwapTotal"]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if merr != nil {
		return nil, merr
	}

	return &m, nil
}
