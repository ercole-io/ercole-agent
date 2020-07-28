// Copyright (c) 2020 Sorint.lab S.p.A.
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

package agentmodel

type ListInstanceOutputModel struct {
	Status      string `json:"status"`
	Name        string `json:"name"`
	ConnString  string `json:"connString"`
	DisplayName string `json:"displayName"`
}

type DbBackupsModel struct {
	DatabaseName string `json:"database_name"`
	Data         []struct {
		BackupType   string  `json:"backup_type"`
		Hour         string  `json:"hour"`
		AvgBckSizeGB float64 `json:"avg_bck_size_gb"`
		WeekDays     string  `json:"week_days"`
	} `json:"data"`
}

type DbSchemasModel struct {
	DatabaseName string `json:"database_name"`
	Data         []struct {
		AllocationType string  `json:"allocation_type"`
		UsedMB         float64 `json:"used_mb"`
		AllocatedMB    float64 `json:"allocated_mb"`
	} `json:"data"`
}

type DbTablespacesModel struct {
	DatabaseName string `json:"database_name"`
	Data         []struct {
		AllocMB    float64 `json:"alloc_mb"`
		UsedMB     float64 `json:"used_mb"`
		Growth     float64 `json:"growth"`
		GrowthUnit string  `json:"growthUnit"`
		FileType   string  `json:"fileType"`
		Filename   string  `json:"file_name"`
		Status     string  `json:"status"`
	} `json:"data"`
}
