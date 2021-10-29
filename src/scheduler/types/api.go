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

package types

import (
	"net/http"
	"time"
)

type Load struct {
	// general info
	SchedulerName string `json:"scheduler_name"`
	// functions-related
	FunctionsDeployed      uint `json:"functions_deployed"`
	FunctionsTotalReplicas uint `json:"functions_total_replicas"`
	FunctionsRunning       uint `json:"functions_running"`
	FunctionsRunningMax    uint `json:"functions_running_max"`
	// queue-related
	QueueLengthMax uint `json:"queue_max_length"`
	QueueFill      int  `json:"queue_fill"`
}

type PeerJobRequest struct {
	// Function    faas.Function     `json:"function"`     // the function that we want to execute
	FunctionName string            `json:"function_name"` // the function name to execute
	Hops         int               `json:"hops"`          // number of times the job is forwarded
	PeersList    []PeersListMember `json:"peers_list"`    // list of peers that handled the job
	Payload      string            `json:"payload"`       // the payload of the request in base64 string
	ContentType  string            `json:"content_type"`  // the mime type of the payload
}

type PeerJobResponse struct {
	PeersList  []PeersListMember `json:"peers_list"`  // list of peers that handled the job
	Body       string            `json:"body"`        // base64 encoded
	StatusCode int               `json:"status_code"` // job response status code
}

type PeersListMember struct {
	MachineId string  `json:"machine_id"`
	MachineIp string  `json:"machine_ip"`
	Timings   Timings `json:"timings"` // timing referred to the passage in that machine
}

type TimingsStart struct {
	ArrivedAt        *time.Time `json:"arrived_at,omitempty"`         // time at which job arrives
	StartedProbingAt *time.Time `json:"started_probing_at,omitempty"` // time at which probing is started
	EndedProbingAt   *time.Time `json:"ended_probing_at,omitempty"`   // time at which probing is ended
	ScheduledAt      *time.Time `json:"scheduled_at,omitempty"`       // time at which job is scheduled internally or externally
}

type Timings struct {
	ExecutionTime  *float64 `json:"execution_time,omitempty"`  // the time of executing the job comprising the GET to openfaas
	TotalTime      *float64 `json:"total_time,omitempty"`      // elapsed time from job arrival til its completed execution
	SchedulingTime *float64 `json:"scheduling_time,omitempty"` // elapsed time for a job to be scheduled
	ProbingTime    *float64 `json:"probing_time,omitempty"`    // elapsed time for a job to probe other nodes
}

type APIResponse struct {
	Headers    http.Header `json:"headers"`
	Body       []byte      `json:"body"`
	StatusCode int         `json:"status_code"`
}
