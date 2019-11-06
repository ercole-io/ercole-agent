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

package model

// Database holds information about the database.
type Database struct {
	InstanceNumber  string
	Name            string
	UniqueName      string
	Status          string
	Version         string
	Platform        string
	Archivelog      string
	Charset         string
	NCharset        string
	BlockSize       string
	CPUCount        string
	SGATarget       string
	PGATarget       string
	MemoryTarget    string
	SGAMaxSize      string
	SegmentsSize    string
	Used            string
	Allocated       string
	Elapsed         string
	DBTime          string
	Work            string
	ASM             bool
	Dataguard       bool
	Patches         []Patch
	Tablespaces     []Tablespace
	Schemas         []Schema
	Features        []Feature
	Licenses        []License
	ADDMs           []Addm
	SegmentAdvisors []SegmentAdvisor
	LastPSUs        []PSU
	Backups         []Backup
	Features2       []Feature2
}
