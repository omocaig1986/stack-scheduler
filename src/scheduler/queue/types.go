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

package queue

import (
	"scheduler/faas"
	"scheduler/types"
	"scheduler/utils"
)

type QueuedJob struct {
	Request   *types.ServiceRequest
	Semaphore *utils.Semaphore
	Response  *faas.APIResponse
	Timings   *Timings
}

type Timings struct {
	ExecutionTime     float64 `json:"execution_time"`      // the time of executing the job comprising the GET to openfaas
	FaasExecutionTime float64 `json:"faas_execution_time"` // the execution time as it is told by openfaas
	QueueTime         float64 `json:"queue_time"`          // the time in which the job remains in the local queue (comprises the execution time)
	// ForwardingTime    float64 `json:"forwarding_time"`     // total time for forwarding the job to another machine
	// ProbingTime       float64 `json:"probing_time"`        // average of time for probing all machines in the fanout (if applicable)
}
