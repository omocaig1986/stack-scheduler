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

package faas

// GetCurrentLoad parse the current load from OpenFaas
func GenGetCurrentLoad(host string) (*CurrentLoad, error) {
	services, _, err := GenFunctionsGet(host)
	if err != nil {
		return nil, err
	}

	replicas := uint(0)
	avReplicas := uint(0)
	nServices := uint(0)
	for _, service := range services {
		replicas += service.Replicas
		avReplicas += service.AvailableReplicas
		nServices += uint(1)
	}

	load := CurrentLoad{
		TotalReplicas:          replicas,
		TotalAvailableReplicas: avReplicas,
		NumberOfServices:       nServices,
	}

	return &load, nil
}

func GenFunctionGetAvailableReplicas(host string, serviceName string) (uint, error) {
	service, _, err := GenFunctionGet(host, serviceName)
	if err != nil {
		return 0, err
	}
	return service.AvailableReplicas, nil
}

func GenFunctionGetReplicas(host string, serviceName string) (uint, error) {
	service, _, err := GenFunctionGet(host, serviceName)
	if err != nil {
		return 0, err
	}
	return service.Replicas, nil
}
