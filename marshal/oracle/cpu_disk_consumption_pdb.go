// Copyright (c) 2023 Sorint.lab S.p.A.
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

package oracle

import (
	"bufio"
	"bytes"
	"strings"
	"time"

	"github.com/ercole-io/ercole-agent/v2/marshal"
	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
	"github.com/hashicorp/go-multierror"
)

// CpuDiskConsumptionsPdb returns information about pluggable database Input / Output Operations Per Second
func CpuDiskConsumptionsPdb(cmdOutput []byte) ([]model.CpuDiskConsumptionPdb, error) {
	res := make([]model.CpuDiskConsumptionPdb, 0)
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	var merr error

	i := 0

	// check if the current line is in the designed marker or not
	isBegin := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "BEGINOUTPUT" {
			isBegin = true
			continue
		}

		if line == "ENDOUTPUT" {
			isBegin = false
			continue
		}

		if !isBegin {
			continue
		}

		splitted := strings.Split(line, "|||")

		// skip row only if all splitted values are N/A
		isValidRow := false
		for _, v := range splitted {
			if !isValidRow {
				isValidRow = v != "N/A"
			}
		}

		if !isValidRow {
			i++
			continue
		}

		if len(splitted) == 5 {
			var start *time.Time

			var end *time.Time

			switch i {
			case 0:
				s := currentTime.AddDate(0, 0, -30)
				start = &s
				end = &currentTime
			case 1:
				s := currentTime.AddDate(0, 0, -7)
				start = &s
				end = &currentTime
			case 2:
				s := currentTime.AddDate(0, 0, -14)
				start = &s
				e := currentTime.AddDate(0, 0, -8)
				end = &e
			case 3:
				s := currentTime.AddDate(0, 0, -21)
				start = &s
				e := currentTime.AddDate(0, 0, -15)
				end = &e
			case 4:
				s := currentTime.AddDate(0, 0, -28)
				start = &s
				e := currentTime.AddDate(0, 0, -22)
				end = &e
			case 5:
				start = &currentTime
				end = nil
			case 6:
				s := currentTime.AddDate(0, 0, -1)
				start = &s
				end = &currentTime
			case 7:
				s := currentTime.AddDate(0, 0, -2)
				start = &s
				e := currentTime.AddDate(0, 0, -1)
				end = &e
			case 8:
				s := currentTime.AddDate(0, 0, -3)
				start = &s
				e := currentTime.AddDate(0, 0, -2)
				end = &e
			case 9:
				s := currentTime.AddDate(0, 0, -4)
				start = &s
				e := currentTime.AddDate(0, 0, -3)
				end = &e
			case 10:
				s := currentTime.AddDate(0, 0, -5)
				start = &s
				e := currentTime.AddDate(0, 0, -4)
				end = &e
			case 11:
				s := currentTime.AddDate(0, 0, -6)
				start = &s
				e := currentTime.AddDate(0, 0, -5)
				end = &e
			}

			sp, err := parseValuesPdb(splitted)
			if err != nil {
				merr = multierror.Append(merr, err)
			}

			sp.TimeStart = start
			sp.TimeEnd = end

			res = append(res, sp)
		}

		if i > 11 {
			sp, err := parseTimeSeriesPdb(splitted)
			if err != nil {
				merr = multierror.Append(merr, err)
			}

			res = append(res, sp)
		}

		i++
	}

	if merr != nil {
		return nil, merr
	}

	return res, nil
}

func parseValuesPdb(lines []string) (model.CpuDiskConsumptionPdb, error) {
	var err, merr error

	sp := model.CpuDiskConsumptionPdb{}

	if sp.CpuDbAvg, err = marshal.TrimParseUnsafeFloat64Pointer(lines[0], marshal.TrimParseFloat64); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if sp.CpuDbMax, err = marshal.TrimParseUnsafeFloat64Pointer(lines[1], marshal.TrimParseFloat64); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if sp.IopsAvg, err = marshal.TrimParseUnsafeFloat64Pointer(lines[2], marshal.TrimParseFloat64); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if sp.IombAvg, err = marshal.TrimParseUnsafeFloat64Pointer(lines[3], marshal.TrimParseFloat64); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if sp.IombMax, err = marshal.TrimParseUnsafeFloat64Pointer(lines[4], marshal.TrimParseFloat64); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	return sp, merr
}

func parseTimeSeriesPdb(lines []string) (model.CpuDiskConsumptionPdb, error) {
	var err, merr error

	sp := model.CpuDiskConsumptionPdb{TimeEnd: nil}

	start, errStart := time.Parse("020115:04", strings.TrimSpace(lines[0]))
	if errStart != nil {
		merr = multierror.Append(merr, ercutils.NewError(errStart))
	}

	start = start.AddDate(time.Now().Year(), 0, 0)
	sp.TimeStart = &start

	if sp.CpuDbAvg, err = marshal.TrimParseUnsafeFloat64Pointer(lines[1], marshal.TrimParseFloat64); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if sp.CpuDbMax, err = marshal.TrimParseUnsafeFloat64Pointer(lines[2], marshal.TrimParseFloat64); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if sp.IopsAvg, err = marshal.TrimParseUnsafeFloat64Pointer(lines[3], marshal.TrimParseFloat64); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if sp.IombAvg, err = marshal.TrimParseUnsafeFloat64Pointer(lines[4], marshal.TrimParseFloat64); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if sp.IombMax, err = marshal.TrimParseUnsafeFloat64Pointer(lines[5], marshal.TrimParseFloat64); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	return sp, merr
}
