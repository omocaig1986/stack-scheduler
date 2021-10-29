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

package scheduler_service

import (
	"scheduler/discovery"
	"scheduler/log"
	"scheduler/utils"
	"sync"
	"time"
)

// getLeastLoadedMachine retrieves the least loaded machine from an array of ips, if all machines are full loaded,
// the least queue is returned, and if there is no less loaded queue than us, an error is returned. This function returns
// (ip, mean_probing_time, errors)
func GetLeastLoadedMachineOfNRandom(n uint, currentLoad uint, checkQueues bool) (string, float64, error) {
	startProbingTime := time.Now()

	// get n random machines from discovery
	machines, err := discovery.GetNRandomMachines(n)
	if err != nil {
		log.Log.Errorf("Cannot get random machines from discovery service")
		return "", 0.0, err
	}

	log.Log.Debugf("len(machines)=%d", len(machines))
	loads := make([]uint, n) // list of loads
	// queues := make([]float64, n) // percentage of queue fill
	probeErr := make([]bool, n) // list of probe errors

	wg := sync.WaitGroup{}
	// get and compute the load of all the available machines in parallel
	for i, ip := range machines {
		wg.Add(1)

		ip := ip
		i := i
		go func() {
			machineLoad, _, err := GetLoad(ip)
			if err != nil {
				log.Log.Errorf("Cannot get load from machine %s", ip)
				probeErr[i] = true
				wg.Done()
				return
			}

			// load := machineLoad.FunctionsRunningMax - machineLoad.FunctionsRunning
			/*
				freeQueue := float64(machineLoad.QueueFill) / float64(machineLoad.QueueLengthMax)
				if freeQueue == math.NaN() {
					log.Log.Debugf("Queue fill value is NaN from machine %s", ip)
					probeErr[i] = true
					wg.Done()
					return
				}
			*/
			loads[i] = uint(machineLoad)
			// queues[i] = 0
			probeErr[i] = false
			wg.Done()
		}()

	}
	wg.Wait()

	log.Log.Debugf("loads=%s", loads)
	log.Log.Debugf("probeErrs=%s", probeErr)

	probingTime := time.Since(startProbingTime).Seconds()

	// Check if we have enough correct loads
	probeErrors := 0
	for i := 0; i < int(n); i++ {
		if probeErr[i] {
			probeErrors += 1
		}
	}
	if probeErrors == int(n) {
		return "", probingTime, NoLessLoadedMachine{"all probe errors"}
	}

	// pick the less loaded
	minLoad, _ := utils.MinOfArrayUint(loads)
	// if no other machine has free slots, see which queue is less loaded
	if minLoad >= currentLoad {
		return "", probingTime, NoLessLoadedMachine{"minLoad >= currentLoad"}
		/*
			==> Queues are no more supported! ==>

			// if we can lose jobs we have not to check queues
			if !checkQueues {
				return "", probingTime, NoLessLoadedMachine{}
			}
			// find the least loaded queue
			leastLoadedQueueValue, leastLoadedQueueIndex := utils.MaxOfArrayFloat(queues)
			ourFreeQueue := float64(queue.GetQueueFill()) / float64(config.Configuration.GetQueueLengthMax())
			if leastLoadedQueueValue < ourFreeQueue {
				return machines[leastLoadedQueueIndex], probingTime, nil
			} else {
				return "", probingTime, NoLessLoadedMachine{}
			}
		*/
	} else {
		// pick one random machine among the less loaded than us
		valuableMachinesIds := utils.LoadsBelowSpecificLoad(loads, currentLoad)
		return machines[valuableMachinesIds[utils.GetRandomInteger(len(valuableMachinesIds))]], probingTime, nil
	}
}
