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

package discovery

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"scheduler/log"
	"scheduler/utils"
	"time"
)

// getListOfServers get the list of known server by asking to the backend stack-service that is running in the same
// machine of this service
func GetListOfMachines() ([]string, error) {
	// get the backend
	res, err := utils.HttpGet(getListApiUrl())
	if err != nil {
		log.Log.Error("Cannot retrieve list of servers from backend")
		return nil, err
	}

	var machines []Machine
	var values []string
	
	response, _ := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()

	err = json.Unmarshal(response, &machines)
	if err != nil {
		log.Log.Error("Cannot unmarshal JSON")
		return nil, err
	}

	for _, machine := range machines {
		values = append(values, machine.IP)
	}

	return values, nil
}

// GetRandomServer returns a N different random servers from the list
func GetNRandomMachines(n uint) ([]string, error) {
	if n == 0 {
		return nil, nil
	}
	list, err := GetListOfMachines()
	if err != nil || len(list) == 0 {
		return nil, &ErrorCannotGetServerList{}
	}
	// if all machines are requested do not pick at random
	if n == uint(len(list)) {
		return list, nil
	}

	randomSource := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(randomSource)

	var out []string
	var picked []int
	for ; ; {
		if len(picked) == int(n) {
			break
		}

		randomI := randomGenerator.Int() % len(list)

		for _, value := range picked {
			if value == randomI {
				continue
			}
		}

		picked = append(picked, randomI)
		out = append(out, list[randomI])
	}
	return out, nil
}

func GetConfiguration() (*ServiceConfiguration, error) {
	res, err := utils.HttpGet(getConfigurationApiUrl())
	if err != nil {
		log.Log.Error("Cannot retrieve configuration from backend")
		return nil, err
	}
	defer res.Body.Close()

	var discoveryConfig ServiceConfiguration
	response, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(response, &discoveryConfig)
	if err != nil {
		log.Log.Error("Cannot unmarshal JSON")
		return nil, err
	}
	return &discoveryConfig, nil
}
