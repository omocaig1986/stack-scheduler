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

package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"scheduler/errors"
	"scheduler/faas"
	"scheduler/log"
	"scheduler/utils"
)

func SystemScaleFunctionPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	function := vars["function"]
	if function == "" {
		errors.ReplyWithError(w, errors.GenericError)
		log.Log.Debugf("service is not specified")
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Log.Debugf("Cannot parse input: %s", err)
		errors.ReplyWithError(w, errors.InputNotValid)
		return
	}

	var params faas.FunctionScalePayload
	err = json.Unmarshal(bytes, &params)
	if err != nil {
		log.Log.Debugf("Cannot parse json input: %s", err)
		errors.ReplyWithError(w, errors.InputNotValid)
		return
	}
	res, err := faas.FunctionScale(function, params.Replicas)
	if err != nil {
		log.Log.Debugf("Cannot scale function: %s", err)
		errors.ReplyWithError(w, errors.GenericOpenFaasError)
		return
	}

	utils.SendJSONResponseByte(&w, res.StatusCode, res.Body)
}
