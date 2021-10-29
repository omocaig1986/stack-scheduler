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
	"encoding/json"
	"io"
	"net/http"
	"scheduler/config"
)

type HelloResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func Hello(w http.ResponseWriter, r *http.Request) {
	helloRes := HelloResponse{
		Name:    config.Name,
		Version: config.Version,
	}
	resBytes, _ := json.Marshal(helloRes)

	_, _ = io.WriteString(w, string(resBytes))
}
