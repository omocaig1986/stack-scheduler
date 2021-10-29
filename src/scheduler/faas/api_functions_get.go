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
	"io/ioutil"
	"scheduler/log"
)

func functionsGetApiCall(host string) (*APIResponse, error) {
	res, err := HttpGet(GetApiSystemFunctionsUrl(host))
	if err != nil {
		log.Log.Errorf("Cannot create GET request to %s", err.Error(), GetApiSystemFunctionsUrl(host))
		return nil, err
	}

	body, _ := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()

	response := APIResponse{
		Headers:    res.Header,
		Body:       body,
		StatusCode: res.StatusCode,
	}

	return &response, err
}
