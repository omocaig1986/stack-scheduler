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
	"scheduler/types"
	"time"
)

const NoSchedulingSchedulerName = "NoScheduler"

type NoSchedulingScheduler struct {
	Loss bool // do not use the queue
}

func (NoSchedulingScheduler) GetFullName() string {
	return NoSchedulingSchedulerName
}

func (s NoSchedulingScheduler) GetScheduler() *types.SchedulerDescriptor {
	return &types.SchedulerDescriptor{
		Name:       NoSchedulingSchedulerName,
		Parameters: []string{fmt.Sprintf("%t", s.Loss)},
	}
}

func (s NoSchedulingScheduler) Schedule(req *types.ServiceRequest) (*JobResult, error) {
	log.Log.Debugf("Scheduling job %s", req.ServiceName)
	now := time.Now()
	timingsStart := types.TimingsStart{ArrivedAt: &now}

	// throw the job if we have no free slots
	if s.Loss && memdb.GetFreeSlots() <= 0 {
		log.Log.Debugf("QueuedJob %s cannot be scheduled, no slots available", req.ServiceName)
		return nil, JobCannotBeScheduled{}
	}

	return executeJobLocally(req, &timingsStart)
}
