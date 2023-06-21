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

var currentTime = time.Now()

// Consumption returns information about host CPU consumption extracted
// from the sar_cpu_only_linux fetcher command output.
func Consumption(cmdOutput []byte) ([]model.Consumption, error) {
	res := []model.Consumption{}
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	var merr, err error

	i := 0

	for scanner.Scan() {
		line := scanner.Text()
		c := model.Consumption{}

		splitted := strings.Split(line, "|||")
		if ercutils.Contains(splitted, "N/A") {
			i++
			continue
		}

		var start *time.Time

		var end *time.Time

		if len(splitted) == 1 {
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

			c.TimeStart = start
			c.TimeEnd = end

			if c.CpuAvg, err = TrimParseUnsafeFloat64(splitted[0], TrimParseFloat64); err != nil {
				merr = multierror.Append(merr, ercutils.NewError(err))
			}
		}

		if len(splitted) == 2 {
			start, errStart := time.Parse("020115:04", strings.TrimSpace(splitted[0]))
			if errStart != nil {
				merr = multierror.Append(merr, ercutils.NewError(errStart))
			}

			start = start.AddDate(time.Now().Year(), 0, 0)
			c.TimeStart = &start
			c.TimeEnd = nil

			if c.CpuAvg, err = TrimParseUnsafeFloat64(splitted[1], TrimParseFloat64); err != nil {
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
