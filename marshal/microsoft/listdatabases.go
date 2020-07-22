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

	"github.com/ercole-io/ercole/model"
)

// ListDatabases marshals -action db output
func ListDatabases(cmdOutput []byte) []model.MicrosoftSQLServerDatabase {
	var rawOut struct {
		Data []struct {
			DatabaseID      int     `json:"database_id"`
			Name            string  `json:"name"`
			CollationName   string  `json:"collation_name"`
			StateDesc       string  `json:"state_desc"`
			RecoveryModel   string  `json:"recovery_model"`
			BlockSize       int     `json:"blocksize"`
			SchedulersCount int     `json:"schedulers_count"`
			AffinityMask    int     `json:"affinity_mask"`
			MinServerMemory int     `json:"min_server_memory"`
			MaxServerMemory int     `json:"max_server_memory"`
			CTP             int     `json:"ctp"`
			MaxDop          int     `json:"maxdop"`
			Alloc           float64 `json:"alloc"`
		} `json:"data"`
	}

	err := json.Unmarshal(cmdOutput, &rawOut)
	if err != nil {
		panic(err)
	}

	out := make([]model.MicrosoftSQLServerDatabase, len(rawOut.Data))

	for i, v := range rawOut.Data {
		out[i].DatabaseID = v.DatabaseID
		out[i].Name = v.Name
		out[i].CollationName = v.CollationName
		out[i].Status = v.StateDesc
		out[i].RecoveryModel = v.RecoveryModel
		out[i].CollationName = v.CollationName
		out[i].BlockSize = v.BlockSize
		out[i].SchedulersCount = v.SchedulersCount
		out[i].AffinityMask = v.AffinityMask
		out[i].MinServerMemory = v.MinServerMemory
		out[i].MaxServerMemory = v.MaxServerMemory
		out[i].CTP = v.CTP
		out[i].MaxDop = v.MaxDop
		out[i].Alloc = v.Alloc
	}

	return out
}
