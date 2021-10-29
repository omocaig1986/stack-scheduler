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
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
)

type IdentifiableFunction struct {
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
}

// ComputeFunctionMD5 computes the ID of a function
func ComputeFunctionMD5(fn *Function) string {
	idFn := IdentifiableFunction{
		Name:         fn.Name,
		Service:      fn.Service,
		Network:      fn.Network,
		Image:        fn.Image,
		EnvProcess:   fn.EnvProcess,
		EnvVars:      fn.EnvVars,
		Constraints:  fn.Constraints,
		Labels:       fn.Labels,
		Annotations:  fn.Annotations,
		Secrets:      fn.Secrets,
		RegistryAuth: fn.RegistryAuth,
		Limits:       fn.Limits,
		Requests:     fn.Requests,
	}

	out, _ := json.Marshal(idFn)
	jsonString := string(out)

	h := md5.New()
	io.WriteString(h, jsonString)
	return fmt.Sprintf("%x", h.Sum(nil))
}
