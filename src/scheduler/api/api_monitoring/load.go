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

package api_monitoring

import (
	"net/http"
	"scheduler/config"
	"scheduler/memdb"
	"strconv"
)

const ApiMonitoringLoadHeaderKey = "X-P2PFaaS-Load"
const ApiMonitoringMaxLoadHeaderKey = "X-P2PFaaS-MaxLoad"

// const ApiMonitoringQueueFillHeaderKey = "X-P2PFog-QueueFill"

// Retrieve the load of the machine.
func LoadGetLoad(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(ApiMonitoringLoadHeaderKey, strconv.Itoa(int(memdb.GetTotalRunningFunctions())))
	w.Header().Add(ApiMonitoringMaxLoadHeaderKey, strconv.Itoa(int(config.Configuration.GetRunningFunctionMax())))
	// w.Header().Add(ApiMonitoringQueueFillHeaderKey, strconv.Itoa(queue.GetQueueFill()))

	w.WriteHeader(200)
	/*
		load, err := faas.GetCurrentLoad()
		if err != nil {
			errors.ReplyWithError(w, errors.GenericError)
			log.Log.Debugf("%s cannot get current load from openfaas")
			return
		}

		res := &types.Load{
			SchedulerName:          scheduler.GetName(),
			FunctionsDeployed:      load.NumberOfServices,
			FunctionsRunning:       memdb.GetTotalRunningFunctions(),
			FunctionsRunningMax:    config.Configuration.GetRunningFunctionMax(),
			FunctionsTotalReplicas: load.TotalReplicas,
			QueueLengthMax:         config.Configuration.GetQueueLengthMax(),
			QueueFill:              queue.GetQueueFill(),
		}

		rep, err := json.Marshal(res)

		utils.SendJSONResponse(&w, 200, string(rep))
	*/
}
