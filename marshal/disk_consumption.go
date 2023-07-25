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

package marshal

import (
	"bufio"
	"bytes"
	"strings"
	"time"

	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
	"github.com/hashicorp/go-multierror"
)

// DiskConsumption returns information about host disk consumption extracted
// from the sar_disks_only_linux fetcher command output.
func DiskConsumption(cmdOutput []byte) ([]model.DiskConsumption, error) {
	res := []model.DiskConsumption{}
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	var merr, err error

	i := 0

	for scanner.Scan() {
		line := scanner.Text()
		c := model.DiskConsumption{}

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

		var start *time.Time

		var end *time.Time

		if len(splitted) == 2 {
			switch i {
			case 0:
				s := currentTime.AddDate(0, 0, -30)
				start = &s
				end = &currentTime
				c.Target = "m"

			case 1:
				s := currentTime.AddDate(0, 0, -7)
				start = &s
				end = &currentTime
				c.Target = "w1"

			case 2:
				s := currentTime.AddDate(0, 0, -14)
				start = &s
				e := currentTime.AddDate(0, 0, -8)
				end = &e
				c.Target = "w2"

			case 3:
				s := currentTime.AddDate(0, 0, -21)
				start = &s
				e := currentTime.AddDate(0, 0, -15)
				end = &e
				c.Target = "w3"

			case 4:
				s := currentTime.AddDate(0, 0, -28)
				start = &s
				e := currentTime.AddDate(0, 0, -22)
				end = &e
				c.Target = "w4"

			case 5:
				start = &currentTime
				end = nil
				c.Target = "d1"

			case 6:
				s := currentTime.AddDate(0, 0, -1)
				start = &s
				end = &currentTime
				c.Target = "d2"

			case 7:
				s := currentTime.AddDate(0, 0, -2)
				start = &s
				e := currentTime.AddDate(0, 0, -1)
				end = &e
				c.Target = "d3"

			case 8:
				s := currentTime.AddDate(0, 0, -3)
				start = &s
				e := currentTime.AddDate(0, 0, -2)
				end = &e
				c.Target = "d4"

			case 9:
				s := currentTime.AddDate(0, 0, -4)
				start = &s
				e := currentTime.AddDate(0, 0, -3)
				end = &e
				c.Target = "d5"

			case 10:
				s := currentTime.AddDate(0, 0, -5)
				start = &s
				e := currentTime.AddDate(0, 0, -4)
				end = &e
				c.Target = "d6"

			case 11:
				s := currentTime.AddDate(0, 0, -6)
				start = &s
				e := currentTime.AddDate(0, 0, -5)
				end = &e
				c.Target = "d7"
			}

			c.TimeStart = start
			c.TimeEnd = end

			if c.IopsHostDayAvg, err = TrimParseUnsafeFloat64Pointer(splitted[0], TrimParseFloat64); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			if c.IombHostDayAvg, err = TrimParseUnsafeFloat64Pointer(splitted[1], TrimParseFloat64); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
		}

		if len(splitted) == 3 {
			start, errStart := time.Parse("020115:04", strings.TrimSpace(splitted[0]))
			if errStart != nil {
				merr = multierror.Append(merr, ercutils.NewError(errStart))
			}

			start = start.AddDate(time.Now().Year(), 0, 0)
			c.TimeStart = &start
			c.TimeEnd = nil

			if c.IopsHostDayAvg, err = TrimParseUnsafeFloat64Pointer(splitted[1], TrimParseFloat64); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}

			if c.IombHostDayAvg, err = TrimParseUnsafeFloat64Pointer(splitted[2], TrimParseFloat64); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
		}

		res = append(res, c)
		i++
	}

	if merr != nil {
		return nil, merr
	}

	return res, nil
}
