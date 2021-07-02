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

package oracle

import (
	"encoding/json"

	"github.com/ercole-io/ercole/v2/model"
	ercutils "github.com/ercole-io/ercole/v2/utils"
)

// Patches marshals -action patches output
func Patches(cmdOutput []byte) ([]model.MicrosoftSQLServerPatch, error) {
	var rawOut struct {
		Data []struct {
			DisplayName    string
			DisplayVersion string
			InstallDate    string
		} `json:"data"`
	}

	if err := json.Unmarshal(cmdOutput, &rawOut); err != nil {
		return nil, ercutils.NewError(err)
	}

	out := make([]model.MicrosoftSQLServerPatch, len(rawOut.Data))

	for i, v := range rawOut.Data {
		out[i].DisplayName = v.DisplayName
		out[i].DisplayVersion = v.DisplayVersion
		out[i].InstallDate = v.InstallDate
	}

	return out, nil
}
