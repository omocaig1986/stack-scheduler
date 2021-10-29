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

import "scheduler/config"

func FunctionsGet() ([]Function, *APIResponse, error) {
	return GenFunctionsGet(config.Configuration.GetOpenFaasListeningHost())
}

func FunctionGet(functionName string) (*Function, *APIResponse, error) {
	return GenFunctionGet(config.Configuration.GetOpenFaasListeningHost(), functionName)
}

func FunctionDeploy(function Function) (*APIResponse, error) {
	return GenFunctionDeploy(config.Configuration.GetOpenFaasListeningHost(), function)
}

func FunctionExecute(functionName string, payload []byte, contentType string) (*APIResponse, error) {
	return GenFunctionExecute(config.Configuration.GetOpenFaasListeningHost(), functionName, payload, contentType)
}

func FunctionScale(functionName string, replicas uint) (*APIResponse, error) {
	return GenFunctionScale(config.Configuration.GetOpenFaasListeningHost(), functionName, replicas)
}

func FunctionScaleByOne(functionName string) (*APIResponse, error) {
	return GenFunctionScaleByOne(config.Configuration.GetOpenFaasListeningHost(), functionName)
}

func FunctionScaleDownByOne(functionName string) (*APIResponse, error) {
	return GenFunctionScaleDownByOne(config.Configuration.GetOpenFaasListeningHost(), functionName)
}
