/*
 * P2PFaaS - A framework for FaaS Load Balancing
 * Copyright (c) 2019. Gabriele Proietti Mattia <pm.gabriele@outlook.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// Package memdb implements a fast way for in-memory variables.
package memdb

import (
	"scheduler/config"
	"scheduler/log"
	"scheduler/metrics"
	"sync"
)

type Function struct {
	Name             string
	RunningInstances uint
}

type ErrorFunctionNotFound struct{}

func (ErrorFunctionNotFound) Error() string {
	return "Function not found"
}

/*
 * Code
 */

var functions []*Function
var totalRunningFunctions uint = 0
var requestNumber uint64 = 0

var mutexRunningFunctions sync.Mutex
var mutexRequestNumber sync.Mutex

func GetRunningInstances(functionName string) (uint, error) {
	mutexRunningFunctions.Lock()
	fn := getFunction(functionName, true)
	if fn == nil {
		mutexRunningFunctions.Unlock()
		return 0, ErrorFunctionNotFound{}
	}
	mutexRunningFunctions.Unlock()
	return fn.RunningInstances, nil
}

func SetFunctionRunning(functionName string) error {
	mutexRunningFunctions.Lock()

	log.Log.Debugf("Setting %s as running", functionName)

	fn := getFunction(functionName, true)
	if fn == nil {
		mutexRunningFunctions.Unlock()
		return ErrorFunctionNotFound{}
	}
	fn.RunningInstances += 1
	totalRunningFunctions += 1
	// metrics
	metrics.PostStartedExecutingJob()

	mutexRunningFunctions.Unlock()
	return nil
}

func SetFunctionStopped(functionName string) error {
	mutexRunningFunctions.Lock()

	log.Log.Debugf("Setting %s as stopped", functionName)

	fn := getFunction(functionName, false)
	if fn == nil {
		mutexRunningFunctions.Unlock()
		return ErrorFunctionNotFound{}
	}
	fn.RunningInstances -= 1
	totalRunningFunctions -= 1
	// metrics
	metrics.PostStoppedExecutingJob()

	mutexRunningFunctions.Unlock()
	return nil
}

func GetTotalRunningFunctions() uint {
	return totalRunningFunctions
}

func GetFreeSlots() int {
	return int(config.Configuration.GetRunningFunctionMax()) - int(totalRunningFunctions)
}

// GetNextRequestNumber returns the next id for the request
func GetNextRequestNumber() uint64 {
	mutexRequestNumber.Lock()
	requestNumber++
	n := requestNumber
	mutexRequestNumber.Unlock()
	return n
}

/*
 * Utils
 */

func getFunction(functionName string, createIfNotExists bool) *Function {
	for _, fn := range functions {
		if fn.Name == functionName {
			return fn
		}
	}

	log.Log.Debugf("%s function not found, creating", functionName)

	if createIfNotExists {
		newFn := Function{
			Name:             functionName,
			RunningInstances: 0,
		}
		functions = append(functions, &newFn)
		return &newFn
	}

	return nil
}
