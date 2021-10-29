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

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	jobHops = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "scheduler_job_hops",
		Help: "Number of hops per job.",
	}, []string{"function_name", "code"})

	jobExecutionTime = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "scheduler_job_execution_time",
		Help: "Time for executing the job in the machine that actually executes it",
	}, []string{"function_name", "code"})

	jobQueueTime = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "scheduler_job_queue_time",
		Help: "Total time for the job to stay in the queue before being executed",
	}, []string{"function_name", "code"})

	jobForwardingTime = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "scheduler_job_forwarding_time",
		Help: "Total time for the job to wait for being executed externally",
	}, []string{"function_name", "code"})

	jobFaasExecutionTime = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "scheduler_job_faas_execution_time",
		Help: "Total time for the job for being executed by openfaas",
	}, []string{"function_name", "code"})

	jobForwardedCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "scheduler_job_forwarded_count",
		Help: "Total time for the job for being executed by openfaas",
	}, []string{"function_name"})

	invocationsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "scheduler_total_invocations",
		Help: "The total number of requests for executing a function",
	}, []string{"function_name", "code"})

	queueFill = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "scheduler_queue_fill",
		Help: "The number of jobs in the queue",
	})

	queueFree = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "scheduler_queue_free",
		Help: "The number of free slots in the queue",
	})

	currentRunningJobs = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "scheduler_current_running_jobs",
		Help: "The number of currently running jobs, not forwarded but executed locally even from remote",
	})

	currentFreeRunningJobs = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "scheduler_current_free_running_jobs",
		Help: "The number of free slots for running jobs, not forwarded but executed locally even from remote",
	})
)
