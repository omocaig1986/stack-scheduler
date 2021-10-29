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

package scheduler

import (
	"fmt"
	"scheduler/log"
	"scheduler/memdb"
	"scheduler/scheduler_service"
	"scheduler/types"
	"time"
)

const PowerOfNSchedulerName = "PowerOfNScheduler"

type PowerOfNScheduler struct {
	F       uint // fan-out
	T       uint // threshold
	Loss    bool // discard job if queue is full
	MaxHops uint // maximum number of hops
}

func (s PowerOfNScheduler) GetFullName() string {
	return fmt.Sprintf("%s(%d, %d, %t, %d)", PowerOfNSchedulerName, s.F, s.T, s.Loss, s.MaxHops)
}

func (s PowerOfNScheduler) GetScheduler() *types.SchedulerDescriptor {
	return &types.SchedulerDescriptor{
		Name: PowerOfNSchedulerName,
		Parameters: []string{
			fmt.Sprintf("%d", s.F),
			fmt.Sprintf("%d", s.T),
			fmt.Sprintf("%t", s.Loss),
			fmt.Sprintf("%d", s.MaxHops),
		},
	}
}

// Schedule a service request. This call is blocking until the job has been executed locally or externally.
func (s PowerOfNScheduler) Schedule(req *types.ServiceRequest) (*JobResult, error) {
	log.Log.Debugf("Scheduling job %s", req.ServiceName)
	currentLoad := memdb.GetTotalRunningFunctions()
	startedScheduling := time.Now()
	timingsStart := types.TimingsStart{ArrivedAt: &startedScheduling}

	balancingHit := currentLoad >= s.T
	jobMustExecutedHere := req.External && req.ExternalJobRequest.Hops >= int(s.MaxHops)

	log.Log.Debugf("balancingHit %t - jobMustExecutedHere %t", balancingHit, jobMustExecutedHere)

	// check if the balancing condition is hit
	if balancingHit && !jobMustExecutedHere {
		// save time
		startedProbingTime := time.Now()
		timingsStart.StartedProbingAt = &startedProbingTime
		// get N Random machines and ask them for load and pick the least loaded
		leastLoaded, _, err := scheduler_service.GetLeastLoadedMachineOfNRandom(s.F, currentLoad, !s.Loss)
		// save time
		endProbingTime := time.Now()
		timingsStart.EndedProbingAt = &endProbingTime

		if err != nil {
			log.Log.Debugf("Error in retrieving machines %s", err.Error())
			// no machine less loaded than us, we are obliged to run the job in this machine or discard the job
			// if we cannot handle it
			return executeJobLocally(req, &timingsStart)
		}

		return executeJobExternally(req, leastLoaded, &timingsStart)
	}

	return executeJobLocally(req, &timingsStart)
}
