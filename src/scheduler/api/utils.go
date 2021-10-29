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

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scheduler/config"
	"scheduler/log"
	"scheduler/scheduler"
	"scheduler/types"
)

const HeaderP2PFaaSVersion = "X-P2PFaaS-Version"
const HeaderP2PFaaSScheduler = "X-P2PFaaS-Scheduler"
const HeaderP2PFaaSTotalTime = "X-P2PFaaS-Timing-Total-Time-Seconds"
const HeaderP2PFaaSExecutionTime = "X-P2PFaaS-Timing-Execution-Time-Seconds"
const HeaderP2PFaaSProbingTime = "X-P2PFaaS-Timing-Probing-Time-Seconds"
const HeaderP2PFaaSSchedulingTime = "X-P2PFaaS-Timing-Scheduling-Time-Seconds"
const HeaderP2PFaaSExternallyExecuted = "X-P2PFaaS-Externally-Executed"
const HeaderP2PFaaSHops = "X-P2PFaaS-Hops"
const HeaderP2PFaaSPeersListIp = "X-P2PFaaS-Peers-List-Ip"
const HeaderP2PFaaSPeersListId = "X-P2PFaaS-Peers-List-Id"

const HeaderP2PFaaSTotalTimingsList = "X-P2PFaaS-Timing-Total-Seconds-List"
const HeaderP2PFaaSProbingTimingsList = "X-P2PFaaS-Timing-Probing-Seconds-List"
const HeaderP2PFaaSSchedulingTimingsList = "X-P2PFaaS-Timing-Scheduling-Seconds-List"

/*
 * Utils
 */

func copyXHeaders(req *types.ServiceRequest, apiResponse *types.APIResponse, toSendResponse *http.ResponseWriter) {
	if apiResponse == nil {
		log.Log.Debugf("[R#%d] apiResponse is nil for job", req.Id)
		return
	}
	for key, value := range apiResponse.Headers {
		if string(key[0]) == "X" {
			log.Log.Debugf("[R#%d] Added header %s = %s", req.Id, key, value[0])
			(*toSendResponse).Header().Set(key, value[0])
		}
	}
}

func addExecuteFunctionCustomHeaders(req *types.ServiceRequest, toSendResponse *http.ResponseWriter, job *scheduler.JobResult) {
	if toSendResponse == nil || job == nil {
		log.Log.Debugf("[R#%d] Cannot add headers: toSend==nil?=%t job==nil?=%t", req.Id, toSendResponse == nil, job == nil)
		return
	}

	(*toSendResponse).Header().Add(HeaderP2PFaaSVersion, config.Version)
	(*toSendResponse).Header().Add(HeaderP2PFaaSScheduler, scheduler.GetName())

	// PowerOfN Headers
	if scheduler.GetName() == scheduler.PowerOfNSchedulerName {
		(*toSendResponse).Header().Add("X-P2PFaaS-Timing-Probe-Messages", fmt.Sprintf("%d", job.ProbingMessages))
	}

	// job has been executed internally so we have single times
	if !job.ExternalExecution {
		if job.Timings == nil {
			log.Log.Errorf("[R#%d] Cannot add timings: they are nil", req.Id)
			return
		}
		if job.Timings.TotalTime != nil {
			(*toSendResponse).Header().Add(HeaderP2PFaaSTotalTime, fmt.Sprintf("%f", *job.Timings.TotalTime))
		}
		if job.Timings.SchedulingTime != nil {
			(*toSendResponse).Header().Add(HeaderP2PFaaSSchedulingTime, fmt.Sprintf("%f", *job.Timings.SchedulingTime))
		}
		if job.Timings.ExecutionTime != nil {
			(*toSendResponse).Header().Add(HeaderP2PFaaSExecutionTime, fmt.Sprintf("%f", *job.Timings.ExecutionTime))
		}
		if job.Timings.ProbingTime != nil {
			(*toSendResponse).Header().Add(HeaderP2PFaaSProbingTime, fmt.Sprintf("%f", *job.Timings.ProbingTime))
		}
	}

	// job has been executed externally so we have a list of times
	if job.ExternalExecution {
		log.Log.Debugf("[R#%d] Job is external and peers list has %d items", req.Id, len(job.ExternalExecutionInfo.PeersList))

		hops := len(job.ExternalExecutionInfo.PeersList) - 1
		(*toSendResponse).Header().Add(HeaderP2PFaaSExternallyExecuted, "True")
		(*toSendResponse).Header().Add(HeaderP2PFaaSHops, fmt.Sprintf("%d", hops))

		if job.ExternalExecutionInfo.PeersList[0].Timings.ExecutionTime != nil {
			(*toSendResponse).Header().Add(HeaderP2PFaaSExecutionTime, fmt.Sprintf("%f", *job.ExternalExecutionInfo.PeersList[0].Timings.ExecutionTime))
		} else {
			log.Log.Errorf("[R#%d] No execution time", req.Id)
		}

		var ipList []string
		var idList []string
		var probingTimes []float64
		var totalTimes []float64
		var schedulingTimes []float64

		peers := job.ExternalExecutionInfo.PeersList
		for i := len(job.ExternalExecutionInfo.PeersList) - 1; i >= 0; i-- {
			ipList = append(ipList, peers[i].MachineIp)
			idList = append(idList, peers[i].MachineId)

			if peers[i].Timings.ProbingTime != nil {
				probingTimes = append(probingTimes, *peers[i].Timings.ProbingTime)
			} else {
				probingTimes = append(probingTimes, 0.0)
			}

			if peers[i].Timings.TotalTime != nil {
				totalTimes = append(totalTimes, *peers[i].Timings.TotalTime)
			} else {
				totalTimes = append(totalTimes, 0.0)
			}

			if peers[i].Timings.SchedulingTime != nil {
				schedulingTimes = append(schedulingTimes, *peers[i].Timings.SchedulingTime)
			} else {
				schedulingTimes = append(schedulingTimes, 0.0)
			}
		}

		ipListJ, _ := json.Marshal(ipList)
		idListJ, _ := json.Marshal(idList)
		totalTimesJ, _ := json.Marshal(totalTimes)
		schedulingTimesJ, _ := json.Marshal(schedulingTimes)
		probingTimesJ, _ := json.Marshal(probingTimes)

		(*toSendResponse).Header().Add(HeaderP2PFaaSPeersListIp, fmt.Sprintf("%s", string(ipListJ)))
		(*toSendResponse).Header().Add(HeaderP2PFaaSPeersListId, fmt.Sprintf("%s", string(idListJ)))
		(*toSendResponse).Header().Add(HeaderP2PFaaSTotalTimingsList, fmt.Sprintf("%s", string(totalTimesJ)))
		(*toSendResponse).Header().Add(HeaderP2PFaaSProbingTimingsList, fmt.Sprintf("%s", string(probingTimesJ)))
		(*toSendResponse).Header().Add(HeaderP2PFaaSSchedulingTimingsList, fmt.Sprintf("%s", string(schedulingTimesJ)))
	}
}
