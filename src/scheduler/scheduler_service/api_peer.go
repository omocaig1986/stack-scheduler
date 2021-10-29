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

package scheduler_service

import (
	"encoding/json"
	"io/ioutil"
	"scheduler/log"
	"scheduler/types"
	"scheduler/utils"
)

func peerFunctionApiCall(host string, request *types.PeerJobRequest) (*APIResponse, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		log.Log.Errorf("Cannot encode to json payload")
		return nil, err
	}

	log.Log.Debugf("Calling POST to %s", GetPeerFunctionUrl(host, request.FunctionName))
	// log.Log.Debugf("len(payload)=%d len(peers)=%d content_type=%s", len(payload), len(request.PeersList), request.ContentType)

	res, err := utils.HttpMachinePostJSON(GetPeerFunctionUrl(host, request.FunctionName), string(payload))
	if err != nil {
		log.Log.Errorf("Cannot create POST request to %s: %s", GetPeerFunctionUrl(host, request.FunctionName), err.Error())
		return nil, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	response := APIResponse{
		Headers:    res.Header,
		Body:       body,
		StatusCode: res.StatusCode,
	}

	return &response, err
}
