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

const PowerOfNSchedulerTauName = "PowerOfNSchedulerTau"

type PowerOfNSchedulerTau struct {
	F       uint          // fan-out
	T       uint          // threshold
	Loss    bool          // discard job if queue is full
	MaxHops uint          // maximum number of hops
	Tau     time.Duration // time to delay the probing
}

func (s PowerOfNSchedulerTau) GetFullName() string {
	return fmt.Sprintf("%s(%d, %d, %t, %d, %dms)", PowerOfNSchedulerTauName, s.F, s.T, s.Loss, s.MaxHops, s.Tau.Milliseconds())
}

func (s PowerOfNSchedulerTau) GetScheduler() *types.SchedulerDescriptor {
	return &types.SchedulerDescriptor{
		Name: PowerOfNSchedulerTauName,
		Parameters: []string{
			fmt.Sprintf("%d", s.F),
			fmt.Sprintf("%d", s.T),
			fmt.Sprintf("%t", s.Loss),
			fmt.Sprintf("%d", s.MaxHops),
			fmt.Sprintf("%dms", s.Tau.Milliseconds()),
		},
	}
}

// Schedule a service request. This call is blocking until the job has been executed locally or externally.
func (s PowerOfNSchedulerTau) Schedule(req *types.ServiceRequest) (*JobResult, error) {
	log.Log.Debugf("[R#%d] Scheduling job %s", req.Id, req.ServiceName)
	currentLoad := memdb.GetTotalRunningFunctions()
	startedScheduling := time.Now()
	timingsStart := types.TimingsStart{ArrivedAt: &startedScheduling}

	balancingHit := currentLoad >= s.T
	jobMustExecutedHere := req.External && req.ExternalJobRequest.Hops >= int(s.MaxHops)

	log.Log.Debugf("[R#%d] balancingHit %t - jobMustExecutedHere %t", req.Id, balancingHit, jobMustExecutedHere)

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
		// compute probing time
		probingTime := endProbingTime.Sub(startedProbingTime)

		// if probing lasted less that Tau wait for reaching tau value
		if time.Since(startedScheduling) < s.Tau {
			time.Sleep(s.Tau - probingTime)
		}

		if err != nil {
			log.Log.Debugf("[R#%d] Error in retrieving machines: %s", req.Id, err.Error())
			// no machine less loaded than us, we are obliged to run the job in this machine or discard the job
			// if we cannot handle it
			return executeJobLocally(req, &timingsStart)
		}

		return executeJobExternally(req, leastLoaded, &timingsStart)
	}

	return executeJobLocally(req, &timingsStart)
}
