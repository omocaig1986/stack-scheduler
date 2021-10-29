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

// Package faas makes possible the communication with the OpenFaaS service.
package faas

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"scheduler/config"
	"scheduler/log"
	"time"
)

var httpTransport *http.Transport

func init() {
	httpTransport = &http.Transport{
		MaxIdleConnsPerHost: 8,
		MaxIdleConns:        24,
		IdleConnTimeout:     64,
		DisableKeepAlives:   false,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 120 * time.Second,
		}).DialContext,
	}
}

/*
 * APIs
 */

func GetApiUrl(host string) string {
	return fmt.Sprintf("http://%s:%d", host, config.Configuration.GetOpenFaasListeningPort())
}

func GetApiSystemFunctionsUrl(host string) string {
	return fmt.Sprintf("%s/system/functions", GetApiUrl(host))
}

func GetApiFunctionUrl(host string, functionName string) string {
	return fmt.Sprintf("%s/function/%s", GetApiUrl(host), functionName)
}

func GetApiSystemFunctionUrl(host string, functionName string) string {
	return fmt.Sprintf("%s/system/function/%s", GetApiUrl(host), functionName)
}

func GetApiScaleFunction(host string, functionName string) string {
	return fmt.Sprintf("%s/system/scale-function/%s", GetApiUrl(host), functionName)
}

/*
 * Http utils
 */

func SetAuthHeader(req *http.Request) {
	auth := config.OpenFaaSUsername + ":" + config.OpenFaaSPassword
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
}

/*
 * Http methods
 */

type ErrorHttpCannotCreateRequest struct{}

func (e ErrorHttpCannotCreateRequest) Error() string {
	return "cannot create http request."
}

func HttpPostJSON(url string, json string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(json))
	if err != nil {
		return nil, ErrorHttpCannotCreateRequest{}
	}

	req.Header.Set("Content-Type", "application/json")
	SetAuthHeader(req)

	client := &http.Client{Transport: httpTransport}
	res, err := client.Do(req)
	if err != nil {
		log.Log.Debugf("cannot POST to %s: %s", url, err.Error())
	}

	return res, err
}

func HttpGet(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, ErrorHttpCannotCreateRequest{}
	}

	SetAuthHeader(req)

	client := &http.Client{Transport: httpTransport}
	res, err := client.Do(req)
	if err != nil {
		log.Log.Debugf("cannot GET to %s: %s", url, err.Error())
	}

	return res, err
}

func HttpPost(url string, payload []byte, contentType string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, ErrorHttpCannotCreateRequest{}
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	SetAuthHeader(req)

	client := &http.Client{Transport: httpTransport}
	res, err := client.Do(req)
	if err != nil {
		log.Log.Debugf("cannot POST to %s: %s", url, err.Error())
	}

	return res, err
}
