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
	"strconv"
)

var executeApiCallResponseHeaderDuration = "X-Duration-Seconds"

func functionExecuteApiCall(host string, functionName string) (*APIResponse, error) {
	res, err := HttpGet(GetApiFunctionUrl(host, functionName))
	if err != nil {
		log.Log.Debugf("Cannot create GET request to %s", err.Error(), GetApiFunctionUrl(host, functionName))
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

func functionExecutePostApiCall(host string, functionName string, payload []byte, contentType string) (*APIResponse, error) {
	res, err := HttpPost(GetApiFunctionUrl(host, functionName), payload, contentType)
	if err != nil {
		log.Log.Debugf("Cannot create POST request to %s", err.Error(), GetApiFunctionUrl(host, functionName))
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

func GetDurationFromExecuteApiCallResponse(res *APIResponse) float64 {
	headerValue := res.Headers.Get(executeApiCallResponseHeaderDuration)
	float, err := strconv.ParseFloat(headerValue, 64)
	if err != nil {
		log.Log.Debugf("Cannot parse float from header value %s: %s", headerValue, err.Error())
		return 0.0
	}
	return float
}
