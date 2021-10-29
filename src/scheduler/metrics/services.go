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

package metrics

import "strconv"

func PostJobMetrics(fnName string, code int, hops int, queueTime float64, execTime float64, faasExecutionTime float64) {
	if enableMetrics {
		jobHops.WithLabelValues(fnName, strconv.Itoa(code)).Observe(float64(hops))
		jobExecutionTime.WithLabelValues(fnName, strconv.Itoa(code)).Observe(execTime)
		jobQueueTime.WithLabelValues(fnName, strconv.Itoa(code)).Observe(queueTime)
		// jobForwardingTime.WithLabelValues(fnName, strconv.Itoa(code)).Observe(forwardingTime)
		jobFaasExecutionTime.WithLabelValues(fnName, strconv.Itoa(code)).Observe(faasExecutionTime)
	}
}

func PostJobInvocations(fnName string, code int) {
	if enableMetrics {
		invocationsTotal.WithLabelValues(fnName, strconv.Itoa(code)).Inc()
	}
}

func PostQueueFreedSlot() {
	if enableMetrics {
		queueFree.Inc()
		queueFill.Dec()
	}
}

func PostQueueAssignedSlot() {
	if enableMetrics {
		queueFree.Dec()
		queueFill.Inc()
	}
}

func PostStartedExecutingJob() {
	if enableMetrics {
		currentRunningJobs.Inc()
		currentFreeRunningJobs.Dec()
	}
}

func PostStoppedExecutingJob() {
	if enableMetrics {
		currentRunningJobs.Dec()
		currentFreeRunningJobs.Inc()
	}
}

func PostJobIsForwarded(fnName string) {
	if enableMetrics {
		jobForwardedCount.WithLabelValues(fnName).Inc()
	}
}

/*
 * Init values set
 */

func PostParallelJobsSlots(n int) {
	if enableMetrics {
		currentFreeRunningJobs.Set(float64(n))
	}
}

func PostQueueSize(n int) {
	if enableMetrics {
		queueFree.Set(float64(n))
	}
}
