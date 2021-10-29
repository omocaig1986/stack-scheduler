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
	"encoding/base64"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"scheduler/config"
	"scheduler/discovery"
	"scheduler/errors"
	"scheduler/log"
	"scheduler/memdb"
	"scheduler/metrics"
	"scheduler/scheduler"
	"scheduler/types"
	"scheduler/utils"
)

func FunctionPost(w http.ResponseWriter, r *http.Request) {
	executeFunction(w, r)
}

func FunctionGet(w http.ResponseWriter, r *http.Request) {
	executeFunction(w, r)
}

/*
 * utils
 */

func executeFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	function := vars["function"]
	if function == "" {
		errors.ReplyWithError(w, errors.GenericError)
		log.Log.Debugf("service is not specified")
		return
	}

	var requestId uint64 = 0
	// assign id to requests if development
	if log.GetEnv() != config.RunningEnvironmentProduction {
		requestId = memdb.GetNextRequestNumber()
	}

	log.Log.Debugf("[R#%d] Execute function called for %s", requestId, function)

	payload, _ := ioutil.ReadAll(r.Body)
	req := types.ServiceRequest{
		Id:          requestId,
		ServiceName: function,
		Payload:     payload,
		ContentType: r.Header.Get("Content-Type"),
		External:    false,
	}

	// schedule the function execution
	jobResult, err := scheduler.Schedule(&req)

	/* This is blocking */

	// check if any error
	if err != nil {
		if cannotScheduleError, ok := err.(*scheduler.JobCannotBeScheduled); ok {
			errors.ReplyWithError(w, errors.JobCannotBeScheduledError)
			log.Log.Debugf("[R#%d] %s", requestId, cannotScheduleError.Error())
			return
		}
		errors.ReplyWithError(w, errors.GenericError)
		log.Log.Debugf("[R#%d] Cannot schedule the service request: %s", requestId, err.Error())
		return
	}

	// check results
	if jobResult != nil && jobResult.Response != nil {
		log.Log.Debugf("[R#%d] Execute function called for %s done: statusCode=%d", requestId, function, jobResult.Response.StatusCode)
	} else if jobResult == nil {
		log.Log.Errorf("[R#%d] jobResult is nil", requestId)
		errors.ReplyWithError(w, errors.GenericError)
		return
	} else if jobResult.Response == nil {
		log.Log.Errorf("[R#%d] jobResult.Response is nil", requestId)
		errors.ReplyWithError(w, errors.GenericError)
		return
	}

	// Compute timings
	utils.ComputeTimings(jobResult.TimingsStart, jobResult.Timings)

	// Add us in list if job is executed externally
	if jobResult.ExternalExecution && jobResult.ExternalExecutionInfo.PeersList != nil {
		jobResult.ExternalExecutionInfo.PeersList = append(
			jobResult.ExternalExecutionInfo.PeersList, discovery.GetPeerDescriptor(jobResult.Timings),
		)
	} else if jobResult.ExternalExecution && jobResult.ExternalExecutionInfo.PeersList == nil {
		log.Log.Errorf("[R#%d] Job has been executed externally but its peers list is empty", requestId)
	}

	// Add custom headers
	copyXHeaders(&req, jobResult.Response, &w)
	addExecuteFunctionCustomHeaders(&req, &w, jobResult)

	w.WriteHeader(jobResult.Response.StatusCode)

	// check if we need to write the body output
	if jobResult.Response.Body != nil && len(jobResult.Response.Body) > 0 {
		log.Log.Debugf("[R#%d] Job body has length %d, external=%t", requestId, len(jobResult.Response.Body), jobResult.ExternalExecution)

		var outputBody []byte

		// decode the job output if it has been executed externally, since when a node offload a jobs, the remote note
		// will reply with base64 encoded body
		if jobResult.ExternalExecution {
			outputBody, err = base64.StdEncoding.DecodeString(string(jobResult.Response.Body))
			if err != nil {
				log.Log.Debugf("[R#%d] Cannot decode job output", requestId)
				return
			}
		} else {
			outputBody = jobResult.Response.Body
		}

		// Write response
		_, err = w.Write(outputBody)
		if err != nil {
			log.Log.Errorf("[R#%d] Cannot write job output: %s", requestId, err.Error())
			return
		}
	}

	// metrics
	defer metrics.PostJobInvocations(function, jobResult.Response.StatusCode)

	defer log.Log.Debugf("[R#%d] %s success", requestId, function)
}
