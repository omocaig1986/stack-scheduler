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
	"net/http"
	"scheduler/errors"
	"scheduler/faas"
	"scheduler/log"
	"scheduler/utils"
)

func SystemFunctionsGet(w http.ResponseWriter, r *http.Request) {
	_, res, err := faas.FunctionsGet()
	if err != nil {
		log.Log.Errorf("Cannot get functions from openfaas")
		errors.ReplyWithError(w, errors.GenericOpenFaasError)
		return
	}

	utils.SendJSONResponseByte(&w, res.StatusCode, res.Body)

	log.Log.Debugf("success")
}

func SystemFunctionsPost(w http.ResponseWriter, r *http.Request) {
	var service faas.Service
	_ = json.NewDecoder(r.Body).Decode(&service)

	res, err := faas.FunctionDeploy(service.OpenFaaSFunction)
	if err != nil {
		errors.ReplyWithError(w, errors.GenericDeployError)
		return
	}

	utils.SendJSONResponseByte(&w, res.StatusCode, res.Body)

	log.Log.Debugf("success")
}

func SystemFunctionsPut(w http.ResponseWriter, r *http.Request) {

}

func SystemFunctionsDelete(w http.ResponseWriter, r *http.Request) {

}
