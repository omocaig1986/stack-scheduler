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
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"scheduler/errors"
	"scheduler/faas"
	"scheduler/log"
	"scheduler/utils"
	"time"
)

type scaleDelayResponse struct {
	ApproximateDelay float64 `json:"approximate_delay"`
	ErrorMargin      float64 `json:"error_margin"`
	Attempts         int     `json:"attempts"`
}

// Measure the time of scaling a service.
func ScaleDelay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	function := vars["function"]

	// before scaling check if the function is stabilized
	fun, _, err := faas.FunctionGet(function)
	if fun.Replicas != fun.AvailableReplicas {
		log.Log.Debugf("Cannot start monitoring, service is not stabilized")
		errors.ReplyWithError(w, errors.InputNotValid)
		return
	}

	// scale the function
	_, err = faas.FunctionScaleByOne(function)
	if err != nil {
		log.Log.Debugf("Cannot scale function: %s", err.Error())
		errors.ReplyWithError(w, errors.GenericOpenFaasError)
		return
	}

	startLoop := time.Now()
	attempts := 0
	var loopTime time.Duration
	// loop until the available replicas is set
	for ; ; {
		fun, _, err := faas.FunctionGet(function)
		attempts += 1
		if err != nil {
			log.Log.Debugf("Cannot get function: %s", err.Error())
			errors.ReplyWithError(w, errors.GenericOpenFaasError)
			break
		}
		if fun.Replicas == fun.AvailableReplicas {
			loopTime = time.Since(startLoop)
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// un-scale the function
	_, err = faas.FunctionScaleDownByOne(function)
	if err != nil {
		log.Log.Debugf("Cannot scale function: %s", err.Error())
		errors.ReplyWithError(w, errors.GenericOpenFaasError)

	}

	res := &scaleDelayResponse{
		Attempts:         attempts,
		ApproximateDelay: loopTime.Seconds(),
		ErrorMargin:      100,
	}
	resJson, _ := json.Marshal(res)
	utils.SendJSONResponse(&w, http.StatusOK, string(resJson))
}
