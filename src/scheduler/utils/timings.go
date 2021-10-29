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
	"scheduler/log"
	"scheduler/types"
	"time"
)

func ComputeTimings(start *types.TimingsStart, timings *types.Timings) {
	log.Log.Debug("Computing timings")
	if start != nil && timings != nil {
		now := time.Now()

		if start.ArrivedAt != nil {
			totalTime := now.Sub(*start.ArrivedAt).Seconds()
			timings.TotalTime = &totalTime
		}

		if start.ScheduledAt != nil {
			schedulingTime := start.ScheduledAt.Sub(*start.ArrivedAt).Seconds()
			timings.SchedulingTime = &schedulingTime
		}

		if start.StartedProbingAt != nil && start.EndedProbingAt != nil {
			probingTime := (start.EndedProbingAt.Sub(*start.StartedProbingAt)).Seconds()
			timings.ProbingTime = &probingTime
		}
	}
}
