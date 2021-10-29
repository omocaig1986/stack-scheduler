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

import "net/http"

type Service struct {
	OpenFaaSFunction Function `json:"openfaas_service,omitempty" bson:"openfaas_service"`
	Deadline         uint64   `json:"deadline,omitempty" bson:"deadline"`
}

type CurrentLoad struct {
	NumberOfServices       uint `json:"total_services" bson:"total_services"`
	TotalReplicas          uint `json:"total_replicas" bson:"total_running_replicas"`
	TotalAvailableReplicas uint `json:"total_available_replicas" bson:"total_available_replicas"`
}

type MachineResources struct {
	Memory string `json:"memory,omitempty" bson:"memory"`
	CPU    string `json:"cpu,omitempty" bson:"cpu"`
}

/*
{
  "service": "nodeinfo",
  "network": "func_functions",
  "image": "functions/nodeinfo:latest",
  "envProcess": "node main.js",
  "envVars": {
    "additionalProp1": "string",
    "additionalProp2": "string",
    "additionalProp3": "string"
  },
  "constraints": [
    "node.platform.os == linux"
  ],
  "labels": [
    "string"
  ],
  "annotations": [
    "string"
  ],
  "secrets": [
    "secret-name-1"
  ],
  "registryAuth": "dXNlcjpwYXNzd29yZA==",
  "limits": {
    "memory": "128M",
    "cpu": "0.01"
  },
  "requests": {
    "memory": "128M",
    "cpu": "0.01"
  }
}
*/

type Function struct {
	Name         string            `json:"name,omitempty" bson:"name"`
	Service      string            `json:"service,omitempty" bson:"service"`
	Network      string            `json:"network,omitempty" bson:"network"`
	Image        string            `json:"image,omitempty" bson:"image"`
	EnvProcess   string            `json:"envProcess,omitempty" bson:"envProcess"`
	EnvVars      map[string]string `json:"envVars,omitempty" bson:"envVars"`
	Constraints  []string          `json:"constraints,omitempty" bson:"constraints"`
	Labels       map[string]string `json:"labels,omitempty" bson:"labels"`
	Annotations  []string          `json:"annotations,omitempty" bson:"annotations"`
	Secrets      []string          `json:"secrets,omitempty" bson:"secrets"`
	RegistryAuth string            `json:"registryAuth,omitempty" bson:"registryAuth"`
	Limits       MachineResources  `json:"limits,omitempty" bson:"limits"`
	Requests     MachineResources  `json:"requests,omitempty" bson:"requests"`

	InvocationCount   uint `json:"invocationCount,omitempty" bson:"invocationCount"`
	Replicas          uint `json:"replicas,omitempty" bson:"replicas"`
	AvailableReplicas uint `json:"availableReplicas,omitempty" bson:"availableReplicas"`
}

type FunctionScalePayload struct {
	Service  string `json:"service,omitempty" bson:"service"`
	Replicas uint   `json:"replicas,omitempty" bson:"replicas"`
}

type APIResponse struct {
	Headers    http.Header
	Body       []byte
	StatusCode int
}
