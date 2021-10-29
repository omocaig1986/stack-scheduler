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
	"github.com/gorilla/mux"
	"net/http"
	"scheduler/errors"
	"scheduler/faas"
	"scheduler/log"
	"scheduler/utils"
)

func SystemFunctionGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	function := vars["function"]
	if function == "" {
		errors.ReplyWithError(w, errors.ServiceNotValid)
		log.Log.Debugf("service is not specified")
		return
	}

	_, res, err := faas.FunctionGet(function)
	if err != nil {
		errors.ReplyWithError(w, errors.GenericOpenFaasError)
		log.Log.Debugf("cannot get the service: %s", err.Error())
		return
	}

	utils.SendJSONResponseByte(&w, res.StatusCode, res.Body)

	log.Log.Debugf("%s success", function)
}
