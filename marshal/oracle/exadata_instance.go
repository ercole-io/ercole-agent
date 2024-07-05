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

func ExadataComponents(cmdOutput []byte) ([]model.OracleExadataComponent, error) {
	components := make([]model.OracleExadataComponent, 0)
	vms := make([]model.OracleExadataVM, 0)
	cells := make([]model.OracleExadataStorageCell, 0)
	grids := make([]model.OracleExadataGridDisk, 0)
	databases := make([]model.OracleExadataDatabase, 0)
	scanner := bufio.NewScanner(bytes.NewReader(cmdOutput))

	var merr error

	for scanner.Scan() {
		var component *model.OracleExadataComponent

		line := scanner.Text()
		splitted := strings.Split(line, "|||")

		if len(splitted) > 0 {
			if splitted[0] == "HOST_TYPE" || splitted[0] == "VM_TYPE" || splitted[0] == "TYPE" {
				continue
			}
		}

		switch splitted[0] {
		case "KVM_HOST", "DOM0", "BARE_METAL":
			if len(splitted) == 18 {
				component, merr = parseKvmHost(splitted)
				components = append(components, *component)
			}
		case "STORAGE_CELL":
			if len(splitted) == 17 {
				component, merr = parseStorageCell(splitted)
				components = append(components, *component)
			}
		case "VM":
			if len(splitted) == 6 {
				components = append(components, *parseVm(splitted))
			}
		case "IB_SWITCH":
			if len(splitted) == 5 {
				components = append(components, *parseIbSwitch(splitted))
			}
		case "VM_KVM":
			if len(splitted) == 8 {
				var vm *model.OracleExadataVM

				vm, merr = parseVmKvm(splitted)
				vms = append(vms, *vm)
			}
		case "VM_XEN":
			if len(splitted) == 7 {
				var vm *model.OracleExadataVM

				vm, merr = parseVmXen(splitted)
				vms = append(vms, *vm)
			}
		case "CELLDISK":
			if len(splitted) == 7 {
				var cell *model.OracleExadataStorageCell

				cell, merr = parseCellDisk(splitted)
				cells = append(cells, *cell)
			}
		case "GRIDDISK":
			if len(splitted) == 11 {
				var grid *model.OracleExadataGridDisk

				grid, merr = parseGridDisk(splitted)
				grids = append(grids, *grid)
			}
		case "DATABASE":
			if len(splitted) == 7 {
				var database *model.OracleExadataDatabase

				database, merr = parseDatabase(splitted)
				databases = append(databases, *database)
			}
		}
	}

	if merr != nil {
		return nil, merr
	}

	return associateExadataVm(
		components,
		vms,
		associateGridToCellDisk(
			grids,
			associateDbToStorageCell(
				databases,
				cells))), nil
}

func associateDbToStorageCell(dbs []model.OracleExadataDatabase, storageCells []model.OracleExadataStorageCell) []model.OracleExadataStorageCell {
	res := make([]model.OracleExadataStorageCell, 0, len(storageCells))

	for _, sc := range storageCells {
		for _, db := range dbs {
			if sc.Cell == db.Cell {
				sc.Databases = append(sc.Databases, db)
			}
		}

		res = append(res, sc)
	}

	return res
}

func associateGridToCellDisk(grids []model.OracleExadataGridDisk, cells []model.OracleExadataStorageCell) []model.OracleExadataStorageCell {
	res := make([]model.OracleExadataStorageCell, 0, len(cells))

	for _, cell := range cells {
		for _, grid := range grids {
			if cell.CellDisk == grid.CellDisk {
				cell.GridDisks = append(cell.GridDisks, grid)
			}
		}

		res = append(res, cell)
	}

	return res
}

func associateExadataVm(components []model.OracleExadataComponent, vms []model.OracleExadataVM, cells []model.OracleExadataStorageCell) []model.OracleExadataComponent {
	res := make([]model.OracleExadataComponent, 0, len(components))

	for _, component := range components {
		for _, vm := range vms {
			if (component.HostType == "KVM_HOST" || component.HostType == "DOM0" || component.HostType == "BARE_METAL") && component.Hostname == vm.PhysicalHost {
				component.VMs = append(component.VMs, vm)
			}
		}

		for _, cell := range cells {
			if component.HostType == "STORAGE_CELL" && component.Hostname == cell.Cell {
				component.StorageCells = append(component.StorageCells, cell)
			}
		}

		res = append(res, component)
	}

	return res
}

func parseKvmHost(sl []string) (*model.OracleExadataComponent, error) {
	res := new(model.OracleExadataComponent)

	var merr, err error

	res.HostType = strings.TrimSpace(sl[0])
	res.RackID = strings.TrimSpace(sl[1])
	res.Hostname = strings.TrimSpace(sl[2])
	res.HostID = strings.TrimSpace(sl[3])

	if res.CPUEnabled, err = marshal.TrimParseInt(sl[4]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.TotalCPU, err = marshal.TrimParseInt(sl[5]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.Memory, err = marshal.TrimParseInt(sl[6]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	res.ImageVersion = strings.TrimSpace(sl[7])
	res.Kernel = strings.TrimSpace(sl[8])
	res.Model = strings.TrimSpace(sl[9])

	if res.FanUsed, err = marshal.TrimParseInt(sl[10]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.FanTotal, err = marshal.TrimParseInt(sl[11]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.PsuUsed, err = marshal.TrimParseInt(sl[12]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.PsuTotal, err = marshal.TrimParseInt(sl[13]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	res.MsStatus = strings.TrimSpace(sl[14])
	res.RsStatus = strings.TrimSpace(sl[15])

	if res.ReservedCPU, err = marshal.TrimParseInt(sl[16]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.ReservedMemory, err = marshal.TrimParseInt(sl[17]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	return res, merr
}

func parseStorageCell(sl []string) (*model.OracleExadataComponent, error) {
	res := new(model.OracleExadataComponent)

	var merr, err error

	res.HostType = strings.TrimSpace(sl[0])
	res.RackID = strings.TrimSpace(sl[1])
	res.Hostname = strings.TrimSpace(sl[2])
	res.HostID = strings.TrimSpace(sl[3])

	if res.CPUEnabled, err = marshal.TrimParseInt(sl[4]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.TotalCPU, err = marshal.TrimParseInt(sl[5]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.Memory, err = marshal.TrimParseInt(sl[6]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	res.ImageVersion = strings.TrimSpace(sl[7])
	res.Kernel = strings.TrimSpace(sl[8])
	res.Model = strings.TrimSpace(sl[9])

	if res.FanUsed, err = marshal.TrimParseInt(sl[10]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.FanTotal, err = marshal.TrimParseInt(sl[11]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.PsuUsed, err = marshal.TrimParseInt(sl[12]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.PsuTotal, err = marshal.TrimParseInt(sl[13]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	res.CellServiceStatus = strings.TrimSpace(sl[14])
	res.MsStatus = strings.TrimSpace(sl[15])
	res.RsStatus = strings.TrimSpace(sl[16])

	return res, merr
}

func parseVm(sl []string) *model.OracleExadataComponent {
	res := new(model.OracleExadataComponent)

	res.HostType = strings.TrimSpace(sl[0])
	res.Hostname = strings.TrimSpace(sl[1])
	res.ImageVersion = strings.TrimSpace(sl[2])
	res.Kernel = strings.TrimSpace(sl[3])
	res.MsStatus = strings.TrimSpace(sl[4])
	res.RsStatus = strings.TrimSpace(sl[5])

	return res
}

func parseIbSwitch(sl []string) *model.OracleExadataComponent {
	res := new(model.OracleExadataComponent)

	res.HostType = strings.TrimSpace(sl[0])
	res.RackID = strings.TrimSpace(sl[1])
	res.Hostname = strings.TrimSpace(sl[2])
	res.Model = strings.TrimSpace(sl[3])
	res.SwVersion = strings.TrimSpace(sl[4])

	return res
}

func parseVmKvm(sl []string) (*model.OracleExadataVM, error) {
	res := new(model.OracleExadataVM)

	var merr, err error

	res.Type = strings.TrimSpace(sl[0])
	res.PhysicalHost = strings.TrimSpace(sl[1])
	res.Status = strings.TrimSpace(sl[2])
	res.Name = strings.TrimSpace(sl[3])

	if res.CPUCurrent, err = marshal.TrimParseInt(sl[4]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.CPURestart, err = marshal.TrimParseInt(sl[5]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.RamCurrent, err = marshal.TrimParseInt(sl[6]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.RamRestart, err = marshal.TrimParseInt(sl[7]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	return res, merr
}

func parseVmXen(sl []string) (*model.OracleExadataVM, error) {
	res := new(model.OracleExadataVM)

	var merr, err error

	res.Type = strings.TrimSpace(sl[0])
	res.PhysicalHost = strings.TrimSpace(sl[1])
	res.Name = strings.TrimSpace(sl[2])

	if res.CPUOnline, err = marshal.TrimParseInt(sl[3]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.CPUMaxUsable, err = marshal.TrimParseInt(sl[4]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.RamOnline, err = marshal.TrimParseInt(sl[5]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.RamMaxUsable, err = marshal.TrimParseInt(sl[6]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	return res, merr
}

func parseCellDisk(sl []string) (*model.OracleExadataStorageCell, error) {
	res := new(model.OracleExadataStorageCell)

	var merr, err error

	res.Type = strings.TrimSpace(sl[0])
	res.CellDisk = strings.TrimSpace(sl[1])
	res.Cell = strings.TrimSpace(sl[2])
	res.Size = strings.TrimSpace(sl[3])
	res.FreeSpace = strings.TrimSpace(sl[4])
	res.Status = strings.TrimSpace(sl[5])

	if res.ErrorCount, err = marshal.TrimParseInt(sl[6]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	return res, merr
}

func parseGridDisk(sl []string) (*model.OracleExadataGridDisk, error) {
	res := new(model.OracleExadataGridDisk)

	var merr, err error

	res.Type = strings.TrimSpace(sl[0])
	res.GridDisk = strings.TrimSpace(sl[1])
	res.CellDisk = strings.TrimSpace(sl[2])
	res.Size = strings.TrimSpace(sl[3])
	res.Status = strings.TrimSpace(sl[4])

	if res.ErrorCount, err = marshal.TrimParseInt(sl[5]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	res.CachingPolicy = strings.TrimSpace(sl[6])
	res.AsmDiskName = strings.TrimSpace(sl[7])
	res.AsmDiskGroup = strings.TrimSpace(sl[8])
	res.AsmDiskSize = strings.TrimSpace(sl[9])
	res.AsmDiskStatus = strings.TrimSpace(sl[10])

	return res, merr
}

func parseDatabase(sl []string) (*model.OracleExadataDatabase, error) {
	res := new(model.OracleExadataDatabase)

	var merr, err error

	res.Type = strings.TrimSpace(sl[0])
	res.DbName = strings.TrimSpace(sl[1])
	res.Cell = strings.TrimSpace(sl[2])

	if res.DbID, err = marshal.TrimParseInt(sl[3]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.FlashCacheLimit, err = marshal.TrimParseInt(sl[4]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if res.IormShare, err = marshal.TrimParseInt(sl[5]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	}

	if t, err := time.Parse(time.RFC3339, sl[6]); err != nil {
		merr = multierror.Append(merr, ercutils.NewError(err))
	} else {
		res.LastIOReq = &t
	}

	return res, merr
}
