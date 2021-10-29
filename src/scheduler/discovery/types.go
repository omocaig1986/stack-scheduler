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

package discovery

type ServiceConfiguration struct {
	MachineIp       string   `json:"machine_ip" bson:"machine_ip"`
	MachineId       string   `json:"machine_id" bson:"machine_id"`
	MachineFogNetId string   `json:"machine_fog_net_id" bson:"machine_fog_net_id"`
	InitServers     []string `json:"init_servers" bson:"init_servers"`
}

type Machine struct {
	ID        int64  `json:"_id" bson:"_id"`
	IP        string `json:"ip" bson:"ip"`
	Name      string `json:"name" bson:"name"`
	GroupName string `json:"group_name" bson:"group_name"`
	// Ping tells the ping, in seconds, of the last poll
	Ping float64 `json:"ping" bson:"ping"` // ms
	// LastUpdate tells the time of the last update
	LastUpdate int64 `json:"last_update" bson:"last_update"`
	// Alive tells if the machine can currently be returned in the list of machine that we known. This parameter is set
	// to false is the machine has been just added or it timed out
	Alive bool `json:"alive" bson:"alive"` // set to not alive when the machine has to be polled
	// DeadPolls tells the number of consecutive times the machine timed out. This is set to 0 when the machine replies
	// correctly
	DeadPolls uint `json:"dead_polls" bson:"dead_polls"`
}
