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

	"github.com/ercole-io/ercole-agent/model"
)

// Features2 returns information about database features2 extracted
// from the opt fetcher command output.
func Features2(cmdOutput []byte) []model.Feature2 {
	features := []model.Feature2{}
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))

	for scanner.Scan() {
		feature := new(model.Feature2)
		line := scanner.Text()
		splitted := strings.Split(line, "|||")
		if len(splitted) == 7 {
			feature.Product = strings.TrimSpace(splitted[0])
			feature.Feature = strings.TrimSpace(splitted[1])
			feature.DetectedUsages = parseInt(strings.TrimSpace(splitted[2]))
			feature.CurrentlyUsed = parseBool(strings.TrimSpace(splitted[3]))
			feature.FirstUsageDate = strings.TrimSpace(splitted[4])
			feature.LastUsageDate = strings.TrimSpace(splitted[5])
			feature.ExtraFeatureInfo = strings.TrimSpace(splitted[6])

			features = append(features, *feature)
		}
	}
	return features
}
