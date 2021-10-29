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

import (
	"encoding/json"
	"io/ioutil"
	"scheduler/log"
)

func functionDeployApiCall(host string, function Function) (*APIResponse, error) {
	faas, err := json.Marshal(function)
	if err != nil {
		log.Log.Debugf("Passed function is not valid: %s", err.Error())
		return nil, err
	}
	log.Log.Debugf("request json is %s", string(faas))

	res, err := HttpPostJSON(GetApiSystemFunctionsUrl(host), string(faas))
	if err != nil {
		log.Log.Debugf("Could not contact OpenFaaS backend: %s", err.Error())
		return nil, err
	}

	body, _ := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()

	response := APIResponse{
		Headers: res.Header,
		Body:       body,
		StatusCode: res.StatusCode,
	}

	return &response, err
}
