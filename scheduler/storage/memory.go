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

// MemoryStorage is a memory task store
type MemoryStorage struct {
	tasks []TaskAttributes
}

// NewMemoryStorage returns an instance of MemoryStorage.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

// Add adds a task to the memory store.
func (memStore *MemoryStorage) Add(task TaskAttributes) error {
	memStore.tasks = append(memStore.tasks, task)
	return nil
}

// Fetch will return all tasks stored.
func (memStore *MemoryStorage) Fetch() ([]TaskAttributes, error) {
	return memStore.tasks, nil
}

// Remove will remove task from store
func (memStore *MemoryStorage) Remove(task TaskAttributes) error {
	var newTasks []TaskAttributes

	for _, existingTask := range memStore.tasks {
		if task.Hash == existingTask.Hash {
			continue
		}

		newTasks = append(newTasks, existingTask)
	}

	memStore.tasks = newTasks

	return nil
}
