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

package api_peer

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"scheduler/config"
	"scheduler/discovery"
	"scheduler/errors"
	"scheduler/log"
	"scheduler/memdb"
	"scheduler/scheduler"
	"scheduler/types"
	"scheduler/utils"
)

// Execute a function. This function must called only by another node, and not a client.
func FunctionExecute(w http.ResponseWriter, r *http.Request) {
	var requestId uint64 = 0
	// assign id to requests if development
	if log.GetEnv() != config.RunningEnvironmentProduction {
		requestId = memdb.GetNextRequestNumber()
	}

	log.Log.Debugf("[R#%d] Request to execute function from peer", requestId)
	vars := mux.Vars(r)
	function := vars["function"]
	if function == "" {
		errors.ReplyWithError(w, errors.GenericError)
		log.Log.Debugf("[R#%d] service is not specified", requestId)
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Log.Debugf("[R#%d] Cannot parse input: %s", requestId, err)
		errors.ReplyWithError(w, errors.InputNotValid)
		return
	}

	var peerRequest types.PeerJobRequest
	err = json.Unmarshal(bytes, &peerRequest)
	if err != nil {
		log.Log.Debugf("[R#%d] Cannot parse json input: %s", requestId, err)
		errors.ReplyWithError(w, errors.InputNotValid)
		return
	}

	req := types.ServiceRequest{
		Id:                 requestId,
		External:           true,
		ExternalJobRequest: &peerRequest,
		ServiceName:        function,
		Payload:            []byte(peerRequest.Payload), // the payload is a string because request it's a peer request
		ContentType:        peerRequest.ContentType,
	}

	log.Log.Debugf("[R#%d] type=%s, len(payload)=%d", requestId, req.ContentType, len(req.Payload))
	log.Log.Debugf("[R#%d] len(peers)=%d, service=%s", requestId, len(peerRequest.PeersList), req.ServiceName)

	// schedule the job
	job, err := scheduler.Schedule(&req)
	// prepare response
	res := preparePeerResponse(&peerRequest, job, err)
	responseBodyBytes, err := json.Marshal(res)

	utils.SendJSONResponse(&w, res.StatusCode, string(responseBodyBytes))
}

// preparePeerResponse Prepares the response to another peer that invoked the function. Remember: jobResult MUST NOT be
// nil even if there is a scheduleErr!
func preparePeerResponse(peerRequest *types.PeerJobRequest, jobResult *scheduler.JobResult, scheduleErr error) *types.PeerJobResponse {
	var res types.PeerJobResponse
	jobBodyResponse := ""

	utils.ComputeTimings(jobResult.TimingsStart, jobResult.Timings)

	// When job ends add us in the peers list
	jobResult.ExternalExecutionInfo.PeersList = append(jobResult.ExternalExecutionInfo.PeersList, discovery.GetPeerDescriptor(jobResult.Timings))

	// Check if result has a response body
	if jobResult.Response != nil {
		// If we have a peer request and we finally executed it here we need to encode the payload in base64
		// We are the last node of the chain PC --> O --> O --> O <-this
		if !jobResult.ExternalExecution {
			// We need to base64 encode the output
			jobBodyResponse = base64.StdEncoding.EncodeToString(jobResult.Response.Body)
		} else {
			jobBodyResponse = string(jobResult.Response.Body) // job result from other nodes is always a base64 string
		}
	}

	if scheduleErr == nil {
		res = types.PeerJobResponse{
			PeersList:  jobResult.ExternalExecutionInfo.PeersList,
			Body:       jobBodyResponse,
			StatusCode: jobResult.Response.StatusCode,
		}
	} else {
		res = types.PeerJobResponse{
			PeersList:  jobResult.ExternalExecutionInfo.PeersList,
			Body:       jobBodyResponse,
			StatusCode: 500,
		}
	}

	return &res
}
