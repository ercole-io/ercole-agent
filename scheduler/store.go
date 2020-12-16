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

package scheduler

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ercole-io/ercole-agent/v2/scheduler/storage"
	"github.com/ercole-io/ercole-agent/v2/scheduler/task"
)

type storeBridge struct {
	store        storage.TaskStore
	funcRegistry *task.FuncRegistry
}

func (sb *storeBridge) Add(task *task.Task) error {
	attributes, err := sb.getTaskAttributes(task)
	if err != nil {
		return err
	}
	return sb.store.Add(attributes)
}

func (sb *storeBridge) Fetch() ([]*task.Task, error) {
	storedTasks, err := sb.store.Fetch()
	if err != nil {
		return []*task.Task{}, err
	}
	var tasks []*task.Task
	for _, storedTask := range storedTasks {
		lastRun, err := time.Parse(time.RFC3339, storedTask.LastRun)
		if err != nil {
			return nil, err
		}

		nextRun, err := time.Parse(time.RFC3339, storedTask.NextRun)
		if err != nil {
			return nil, err
		}

		duration, err := time.ParseDuration(storedTask.Duration)
		if err != nil {
			return nil, err
		}

		isRecurring, err := strconv.Atoi(storedTask.IsRecurring)
		if err != nil {
			return nil, err
		}

		funcMeta, err := sb.funcRegistry.Get(storedTask.Name)
		if err != nil {
			return nil, err
		}

		params, err := paramsFromString(funcMeta, storedTask.Params)
		if err != nil {
			return nil, err
		}

		t := task.NewWithSchedule(funcMeta, params, task.Schedule{
			IsRecurring: isRecurring == 1,
			Duration:    time.Duration(duration),
			LastRun:     lastRun,
			NextRun:     nextRun,
		})
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (sb *storeBridge) Remove(task *task.Task) error {
	attributes, err := sb.getTaskAttributes(task)
	if err != nil {
		return err
	}
	return sb.store.Remove(attributes)
}

func (sb *storeBridge) getTaskAttributes(task *task.Task) (storage.TaskAttributes, error) {
	params, err := paramsToString(task.Params)
	if err != nil {
		return storage.TaskAttributes{}, err
	}

	isRecurring := 0
	if task.IsRecurring {
		isRecurring = 1
	}

	return storage.TaskAttributes{
		Hash:        string(task.Hash()),
		Name:        task.Func.Name,
		LastRun:     task.LastRun.Format(time.RFC3339),
		NextRun:     task.NextRun.Format(time.RFC3339),
		Duration:    task.Duration.String(),
		IsRecurring: strconv.Itoa(isRecurring),
		Params:      params,
	}, nil
}

func paramsToString(params []task.Param) (string, error) {
	var paramsList []string
	for _, param := range params {
		paramStr, err := json.Marshal(param)
		if err != nil {
			return "", err
		}
		paramsList = append(paramsList, string(paramStr))
	}
	data, err := json.Marshal(paramsList)
	return string(data), err
}

func paramsFromString(funcMeta task.FunctionMeta, payload string) ([]task.Param, error) {
	var params []task.Param
	if strings.TrimSpace(payload) == "" {
		return params, nil
	}
	paramTypes := funcMeta.Params()
	var paramsStrings []string
	err := json.Unmarshal([]byte(payload), &paramsStrings)
	if err != nil {
		return params, err
	}
	for i, paramStr := range paramsStrings {
		paramType := paramTypes[i]
		target := reflect.New(paramType)
		err := json.Unmarshal([]byte(paramStr), target.Interface())
		if err != nil {
			return params, err
		}
		param := reflect.Indirect(target).Interface().(task.Param)
		params = append(params, param)
	}

	return params, nil
}
