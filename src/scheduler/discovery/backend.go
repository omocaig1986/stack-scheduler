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
	"fmt"
	"scheduler/config"
)

func getApiUrl() string {
	return fmt.Sprintf("http://%s:%d", config.Configuration.GetDiscoveryListeningHost(), config.Configuration.GetDiscoveryListeningPort())
}

// getListApiUrl returns the api url of the list function for the discovery service
func getListApiUrl() string {
	return fmt.Sprintf("%s/list", getApiUrl())
}

// getListApiUrl returns the api url of the list function for the discovery service
func getConfigurationApiUrl() string {
	return fmt.Sprintf("%s/configuration", getApiUrl())
}
