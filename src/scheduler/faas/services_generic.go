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
	"scheduler/log"
)

func GenFunctionsGet(host string) ([]Function, *APIResponse, error) {
	res, err := functionsGetApiCall(host)
	if err != nil {
		log.Log.Debugf("Cannot get services from openfaas: %s", err.Error())
		return nil, res, err
	}

	var functions []Function
	err = json.Unmarshal([]byte(res.Body), &functions)
	if err != nil {
		log.Log.Debugf("Cannot decode services from openfaas: %s: %s", res.Body, err.Error())
		return nil, res, err
	}
	return functions, res, nil
}

func GenFunctionGet(host string, functionName string) (*Function, *APIResponse, error) {
	res, err := functionGetApiCall(host, functionName)
	if err != nil {
		log.Log.Debugf("Cannot get service from openfaas: %s", err.Error())
		return nil, res, err
	}

	var function Function
	err = json.Unmarshal([]byte(res.Body), &function)

	if err != nil {
		log.Log.Debugf("Reply is (%d) [%s]: %s", res.StatusCode, res.StatusCode, res.Body)
		log.Log.Debugf("Cannot decode service from openfaas: %s", err.Error())
		return nil, res, err
	}
	return &function, res, nil
}

func GenFunctionDeploy(host string, function Function) (*APIResponse, error) {
	res, err := functionDeployApiCall(host, function)
	if err != nil {
		log.Log.Debugf("Cannot deploy service" + function.Name + ": " + err.Error())
	}
	return res, err
}

func GenFunctionExecute(host string, functionName string, payload []byte, contentType string) (*APIResponse, error) {
	var res *APIResponse
	var err error

	if payload == nil {
		res, err = functionExecuteApiCall(host, functionName)
	} else {
		res, err = functionExecutePostApiCall(host, functionName, payload, contentType)
	}

	if err != nil {
		return nil, err
	}
	if res.StatusCode == 404 {
		return res, ErrorFunctionNotFound{}
	}
	if res.StatusCode >= 500 {
		return res, ErrorInternal{string(res.Body)}
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return res, ErrorGeneric{string(res.Body)}
	}
	return res, nil
}

func GenFunctionScale(host string, functionName string, replicas uint) (*APIResponse, error) {
	res, err := functionScaleApiCall(host, functionName, replicas)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 404 {
		return res, ErrorFunctionNotFound{}
	}
	if res.StatusCode >= 500 {
		return res, ErrorInternal{string(res.Body)}
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return res, ErrorGeneric{string(res.Body)}
	}
	return res, nil
}

func GenFunctionScaleByOne(host string, functionName string) (*APIResponse, error) {
	reps, err := GenFunctionGetReplicas(host, functionName)
	if err != nil {
		log.Log.Debugf("Could not scale by one service %s: %s", functionName, err.Error())
		return nil, err
	}
	return functionScaleApiCall(host, functionName, reps+1)
}

func GenFunctionScaleDownByOne(host string, functionName string) (*APIResponse, error) {
	reps, err := GenFunctionGetReplicas(host, functionName)
	if err != nil {
		log.Log.Debugf("Could not scale down by one service %s: %s", functionName, err.Error())
		return nil, err
	}
	return functionScaleApiCall(host, functionName, reps-1)
}
