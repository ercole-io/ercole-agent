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
	"time"

	"github.com/ercole-io/ercole/model"
)

// DatabaseFeatureUsageStat returns information about database features2 extracted
// from the opt fetcher command output.
func DatabaseFeatureUsageStat(cmdOutput []byte) []model.OracleDatabaseFeatureUsageStat {
	featuresUsageStats := []model.OracleDatabaseFeatureUsageStat{}
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))

	for scanner.Scan() {
		stats := new(model.OracleDatabaseFeatureUsageStat)
		line := scanner.Text()
		splitted := strings.Split(line, "|||")

		if len(splitted) == 7 {
			stats.Product = strings.TrimSpace(splitted[0])
			stats.Feature = strings.TrimSpace(splitted[1])
			stats.DetectedUsages = trimParseInt64(splitted[2])
			stats.CurrentlyUsed = parseBool(strings.TrimSpace(splitted[3]))

			var err error
			stats.FirstUsageDate, err = time.Parse("2006-01-02 15:04:05", strings.TrimSpace(splitted[4]))
			if err != nil {
				panic(err)
			}

			stats.LastUsageDate, _ = time.Parse("2006-01-02 15:04:05", strings.TrimSpace(splitted[5]))
			if err != nil {
				panic(err)
			}

			stats.ExtraFeatureInfo = strings.TrimSpace(splitted[6])

			featuresUsageStats = append(featuresUsageStats, *stats)
		}
	}
	return featuresUsageStats
}
