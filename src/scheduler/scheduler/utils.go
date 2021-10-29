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

package scheduler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"scheduler/log"
	"scheduler/memdb"
	"scheduler/queue"
	"scheduler/scheduler_service"
	"scheduler/types"
	"time"
)

/*
 * Core
 */

func executeJobExternally(req *types.ServiceRequest, remoteNodeIP string, timingsStart *types.TimingsStart) (*JobResult, error) {
	log.Log.Debugf("[R#%d] %s scheduled to be run at %s", req.Id, req.ServiceName, remoteNodeIP)
	now := time.Now()
	if timingsStart != nil {
		timingsStart.ScheduledAt = &now
	}

	// prepare everything to send the job externally
	peerRequest, err := prepareForwardToPeerRequest(req)
	if err != nil {
		return nil, err
	}
	// metrics
	// metrics.PostJobIsForwarded(req.ServiceName)

	res, err := scheduler_service.ExecuteFunction(remoteNodeIP, peerRequest)
	/* This is blocking */

	return prepareJobResultFromExternalExecution(req, res, timingsStart), nil
}

func executeJobLocally(req *types.ServiceRequest, timingsStart *types.TimingsStart) (*JobResult, error) {
	log.Log.Debugf("[R#%d] %s scheduled to be run locally: external=%t", req.Id, req.ServiceName, req.External)
	now := time.Now()
	if timingsStart != nil {
		timingsStart.ScheduledAt = &now
	}

	freeSlots := memdb.GetFreeSlots()
	if memdb.GetFreeSlots() <= 0 {
		log.Log.Debugf("[R#%d] %s cannot be scheduled to be run locally: freeSlots=%d", req.Id, req.ServiceName, freeSlots)
		return &JobResult{
			Response:          nil,
			Timings:           &types.Timings{},
			TimingsStart:      timingsStart,
			ExternalExecution: false,
		}, JobCannotBeScheduled{}
	}

	// If we execute the job locally and request is external, payload is base64encoded and we decode it
	if req.External {
		decodedPayload, _ := base64.StdEncoding.DecodeString(string(req.Payload))
		req.Payload = decodedPayload
	}

	job, err := queue.EnqueueJob(req)

	/* This is blocking */

	if err != nil {
		log.Log.Debugf("[R#%d] Cannot add job to queue, job is discarded", req.Id)
		return nil, err
	}

	// Fill the execution time since it is derived from the internal execution
	executionTime := job.Timings.FaasExecutionTime
	timings := types.Timings{ExecutionTime: &executionTime}

	return prepareJobResultFromInternalExecution(job, req, timingsStart, &timings), nil
}

// prepareJobResultFromInternalExecution prepare the result when the job is execute internally
func prepareJobResultFromInternalExecution(job *queue.QueuedJob, req *types.ServiceRequest, timingsStart *types.TimingsStart, timings *types.Timings) *JobResult {
	log.Log.Debugf("[R#%d] status_code=%d", req.Id, job.Response.StatusCode)

	response := types.APIResponse{
		Headers:    job.Response.Headers,
		StatusCode: job.Response.StatusCode,
		Body:       job.Response.Body,
	}

	result := JobResult{
		Response:          &response,
		Timings:           timings,
		TimingsStart:      timingsStart,
		ExternalExecution: false,
	}
	return &result
}

func prepareJobResultFromExternalExecution(req *types.ServiceRequest, res *scheduler_service.APIResponse, timingsStart *types.TimingsStart) *JobResult {
	var response types.APIResponse
	var result JobResult

	// Response should be never nil but we check here in case
	if res != nil {
		log.Log.Debugf("[R#%d] Response from external execution is %d", req.Id, res.StatusCode)
		response = types.APIResponse{
			Headers:    res.Headers,
			StatusCode: res.StatusCode,
			Body:       res.Body,
		}

		var peerJobResponse types.PeerJobResponse
		err := json.Unmarshal(res.Body, &peerJobResponse)
		if err != nil {
			log.Log.Debugf("Cannot decode the job response")
		}

		// Change the body of the response leaving only the output of the function, this because when executing functions
		// between peer nodes we encapsulate the output body in a PeerJobResponse struct
		response.Body = []byte(peerJobResponse.Body)

		// Prepare the result
		result = JobResult{
			Response: &response,
			ExternalExecutionInfo: ExternalExecutionInfo{
				PeersList: peerJobResponse.PeersList,
			},
		}
	} else {
		log.Log.Errorf("[R#%d] Response from peer is null", req.Id)
		response = types.APIResponse{
			Headers:    http.Header{},
			StatusCode: 500,
			Body:       []byte{},
		}

		// prepare the result
		result = JobResult{
			Response:              &response,
			ExternalExecutionInfo: ExternalExecutionInfo{},
		}
	}

	result.TimingsStart = timingsStart
	result.ExternalExecution = true
	result.Timings = &types.Timings{}

	return &result
}

/*
 * Utils
 */

// prepareForwardToPeerRequest prepare the request to execute the job to another peer
func prepareForwardToPeerRequest(req *types.ServiceRequest) (*types.PeerJobRequest, error) {
	peerRequest := types.PeerJobRequest{
		FunctionName: req.ServiceName,
		ContentType:  req.ContentType,
	}

	// If request is external the payload is already in base64
	if !req.External {
		// encode payload in base64
		peerRequest.Payload = base64.StdEncoding.EncodeToString(req.Payload)
		peerRequest.Hops += 1
	} else {
		peerRequest.Payload = string(req.Payload)
		peerRequest.Hops = 1
	}
	return &peerRequest, nil
}
