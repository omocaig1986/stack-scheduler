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
	"sync"
	"time"
)

const RoundRobinWithMasterSchedulerName = "RoundRobinWithMasterScheduler"

type RoundRobinWithMasterScheduler struct {
	Master            bool       // if current node is master
	MasterIP          string     // ip of master node
	Loss              bool       // discard job if queue is full
	currentIndex      int        // current index of the round robin
	currentIndexMutex sync.Mutex // protect race conditions on currentIndex
}

func (s *RoundRobinWithMasterScheduler) GetFullName() string {
	return fmt.Sprintf("%s(%t, %s, %t)", RoundRobinWithMasterSchedulerName, s.Master, s.MasterIP, s.Loss)
}

func (s *RoundRobinWithMasterScheduler) GetScheduler() *types.SchedulerDescriptor {
	return &types.SchedulerDescriptor{
		Name: RoundRobinWithMasterSchedulerName,
		Parameters: []string{
			fmt.Sprintf("%t", s.Master),
			fmt.Sprintf("%s", s.MasterIP),
			fmt.Sprintf("%t", s.Loss),
		},
	}
}

func (s *RoundRobinWithMasterScheduler) Schedule(req *types.ServiceRequest) (*JobResult, error) {
	log.Log.Debugf("Scheduling job %s", req.ServiceName)
	now := time.Now()
	timingsStart := types.TimingsStart{ArrivedAt: &now}

	// Master node
	if s.Master {
		// Master node cannot schedule jobs only dispatch them
		if !req.External {
			return nil, JobCannotBeScheduled{}
		}
		// Obtain the list of all machines and select one with a round robin fashion
		machinesIp, err := discovery.GetListOfMachines()
		if err != nil {
			return nil, JobCannotBeScheduled{err.Error()}
		}
		if len(machinesIp) == 0 {
			return nil, JobCannotBeScheduled{"no machine known"}
		}

		// Update the id of next machine
		s.currentIndexMutex.Lock()
		// Check if current index is not exceeding the length of machines array
		if s.currentIndex >= len(machinesIp) {
			s.currentIndex = 0
		}
		pickedMachineIp := machinesIp[s.currentIndex]
		s.currentIndex = (s.currentIndex + 1) % len(machinesIp)
		s.currentIndexMutex.Unlock()

		log.Log.Debugf("nextIndex is %d", s.currentIndex)

		// Schedule the job to that machine
		return executeJobExternally(req, pickedMachineIp, &timingsStart)
	}

	// Slave node
	// If request is internal dispatch it to the master node
	if !req.External {
		return executeJobExternally(req, s.MasterIP, &timingsStart)
	}

	// Otherwise execute it internally
	return executeJobLocally(req, &timingsStart)
}
