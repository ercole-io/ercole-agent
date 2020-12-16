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
)

// LicensingInfo marshals -action licensinginfo output
func LicensingInfo(cmdOutput []byte, inst *model.MicrosoftSQLServerInstance) {
	var out struct {
		Data struct {
			ProductVersion string
			EditionType    string
			ProductCode    string
			LicensingInfo  string
		} `json:"data"`
	}

	err := json.Unmarshal(cmdOutput, &out)
	if err != nil {
		panic(err)
	}

	inst.Version = out.Data.ProductVersion
	inst.EditionType = out.Data.EditionType
	inst.ProductCode = out.Data.ProductCode
	inst.LicensingInfo = out.Data.LicensingInfo
}
