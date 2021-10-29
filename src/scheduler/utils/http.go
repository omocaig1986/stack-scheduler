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

package utils

import (
	"bytes"
	"io"
	"net/http"
	"scheduler/log"
)

type ErrorHttpCannotCreateRequest struct{}

func (e ErrorHttpCannotCreateRequest) Error() string {
	return "cannot create http request."
}

/*
 * Generic Http methods
 */

func HttpPostJSON(url string, json string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(json))
	if req == nil {
		return nil, ErrorHttpCannotCreateRequest{}
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Transport: httpTransport}
	res, err := client.Do(req)
	if err != nil {
		log.Log.Debugf("Cannot POST to %s: %s", url, err.Error())
	}

	return res, err
}

func HttpGet(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if req == nil {
		return nil, ErrorHttpCannotCreateRequest{}
	}

	client := &http.Client{Transport: httpTransport}
	res, err := client.Do(req)
	if err != nil {
		log.Log.Debugf("Cannot GET to %s: %s", url, err.Error())
	}

	return res, err
}

// HttpMachineGet performs and http get setting as user agent Machine
func HttpMachineGet(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if req == nil {
		return nil, ErrorHttpCannotCreateRequest{}
	}

	req.Header.Add("User-Agent", "Machine")

	client := &http.Client{Transport: httpTransport}
	res, err := client.Do(req)
	if err != nil {
		log.Log.Debugf("Cannot GET to %s: %s", url, err.Error())
	}

	return res, err
}

func HttpMachinePostJSON(url string, json string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(json))
	if req == nil {
		return nil, ErrorHttpCannotCreateRequest{}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Machine")

	client := &http.Client{Transport: httpTransport}
	res, err := client.Do(req)
	if err != nil {
		log.Log.Debugf("Cannot POST to %s: %s", url, err.Error())
	}

	return res, err
}

/*
* Utils
 */

func SendJSONResponse(w *http.ResponseWriter, code int, body string) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(code)

	_, err := io.WriteString(*w, body)
	if err != nil {
		log.Log.Debugf("Cannot send response: %s", err.Error())
	}
}

func SendJSONResponseByte(w *http.ResponseWriter, code int, body []byte) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(code)

	_, err := (*w).Write(body)
	if err != nil {
		log.Log.Debugf("Cannot send response: %s", err.Error())
	}
}
