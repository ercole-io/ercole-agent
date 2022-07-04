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

package postgresql

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"

	"github.com/ercole-io/ercole-agent/v2/marshal"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/hashicorp/go-multierror"
)

func Setting(cmdOutput []byte) (*model.PostgreSQLSetting, error) {
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))
	result := model.PostgreSQLSetting{}

	var merr, err error

	for scanner.Scan() {
		line := scanner.Text()

		splitted := strings.Split(line, "|")
		if len(splitted) == 0 {
			continue
		}

		iter := marshal.NewIter(splitted)

		if len(splitted) == 16 {
			result.DbVersion = iter()

			if result.WorkMem, err = strconv.Atoi(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			result.ArchiveMode = marshal.TrimParseBool(iter())

			result.ArchiveCommand = iter()

			if result.MinWalSize, err = strconv.Atoi(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			if result.MaxWalSize, err = strconv.Atoi(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			if result.MaxConnections, err = strconv.Atoi(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			result.CheckpointCompletionTarget = iter()

			if result.DefaultStatisticsTarget, err = strconv.Atoi(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			if result.RandomPageCost, err = strconv.ParseFloat(iter(), 64); err != nil {
				merr = multierror.Append(merr, err)
			}

			if result.MaintenanceWorkMem, err = strconv.Atoi(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			if result.SharedBuffers, err = strconv.Atoi(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			if result.EffectiveCacheSize, err = strconv.Atoi(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			if result.EffectiveIoConcurrency, err = strconv.Atoi(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			if result.MaxWorkerProcesses, err = strconv.Atoi(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}

			if result.MaxParallelWorkers, err = strconv.Atoi(iter()); err != nil {
				merr = multierror.Append(merr, err)
			}
		}
	}

	if merr != nil {
		return nil, merr
	}

	return &result, nil
}
