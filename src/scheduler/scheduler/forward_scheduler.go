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
	"scheduler/discovery"
	"scheduler/log"
	"scheduler/types"
	"time"
)

const ForwardSchedulerName = "ForwardScheduler"

// This scheduler forward all of its jobs to a random node, this is used for testing purposes
type ForwardScheduler struct {
	MaxHops uint // maximum number of hops
}

func (s ForwardScheduler) GetFullName() string {
	return fmt.Sprintf("%s(%d)", ForwardSchedulerName, s.MaxHops)
}

func (s ForwardScheduler) GetScheduler() *types.SchedulerDescriptor {
	return &types.SchedulerDescriptor{
		Name: ForwardSchedulerName,
		Parameters: []string{
			fmt.Sprintf("%d", s.MaxHops),
		},
	}
}

// Schedule a service request. This call is blocking until the job has been executed locally or externally.
func (s ForwardScheduler) Schedule(req *types.ServiceRequest) (*JobResult, error) {
	log.Log.Debugf("Scheduling job %s", req.ServiceName)
	now := time.Now()
	timingsStart := types.TimingsStart{ArrivedAt: &now}

	jobMustExecutedHere := req.External && req.ExternalJobRequest.Hops >= int(s.MaxHops)

	// check if the balancing condition is hit
	if !jobMustExecutedHere {
		// save time
		startedProbingTime := time.Now()
		timingsStart.StartedProbingAt = &startedProbingTime
		// get N Random machines and ask them for load and pick the least loaded
		randomMachine, err := discovery.GetNRandomMachines(1)
		// save time
		endProbingTime := time.Now()
		timingsStart.EndedProbingAt = &endProbingTime
		if err != nil {
			log.Log.Debugf("Error in retrieving machines %s", err.Error())
			return executeJobLocally(req, &timingsStart)
		}
		if len(randomMachine) == 0 {
			log.Log.Debugf("No random machines retrieved")
			return executeJobLocally(req, &timingsStart)
		}
		log.Log.Debugf("Forwarding to random machine: %s", randomMachine)
		return executeJobExternally(req, randomMachine[0], &timingsStart)
	}

	return executeJobLocally(req, &timingsStart)
}
