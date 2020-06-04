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

package storage

// TaskAttributes is a struct which is used to transfer data from/to stores.
// All task data are converted from/to string to prevent the store from
// worrying about details of converting data to the proper formats.
type TaskAttributes struct {
	Hash        string
	Name        string
	LastRun     string
	NextRun     string
	Duration    string
	IsRecurring string
	Params      string
}

// TaskStore is the interface to implement when adding custom task storage.
type TaskStore interface {
	Add(TaskAttributes) error
	Fetch() ([]TaskAttributes, error)
	Remove(TaskAttributes) error
}
