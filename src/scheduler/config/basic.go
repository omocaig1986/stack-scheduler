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

package config

import (
	"io/ioutil"
	"scheduler/log"
	"strings"
)

const Name = "p2pfaas-scheduler"
const Version = "0.2.3b"

const DataPath = "/data"

// const ConfigurationFilePath = "/config"
const ConfigurationFileName = "p2p_faas-scheduler.json"
const ConfigurationSchedulerFileName = "p2p_faas-scheduler-config.json"

// const ConfigurationFileFullPath = ConfigurationFilePath + "/" + ConfigurationFileName
// const SchedulerConfigurationFullPath = ConfigurationFilePath + "/" + SchedulerConfigurationFileName

const DefaultListeningPort = 18080
const DefaultQueueLengthMax = 100
const DefaultFunctionsRunningMax = 10

// env
const EnvRunningEnvironment = "P2PFAAS_DEV_ENV"
const EnvProfiling = "P2PFAAS_PROF"

// stack-discovery
const DefaultStackDiscoveryListeningHost = "discovery"
const DefaultStackDiscoveryListeningPort = 19000

// environments
const RunningEnvironmentProduction = "production"
const RunningEnvironmentDevelopment = "development"

// OpenFaaS
const DefaultOpenFaaSListeningHost = "faas-swarm"
const DefaultOpenFaaSListeningPort = 8080

/*
 * Variables
 */
var OpenFaaSUsername = "admin"
var OpenFaaSPassword = "admin"

var Configuration ConfigurationSet
var ConfigurationReadFromFile = false

func init() {
	// read the config file
	_, err := ReadConfigFile()
	if err != nil {
		log.Log.Warningf("Cannot read config file at %s, using default values", GetConfigFilePath())
	} else {
		ConfigurationReadFromFile = true
	}

	// get the secrets for accessing OpenFaas APIs
	// if os.Getenv(EnvDevelopmentEnvironment) == Configuration.GetRunningEnvironment() {
	log.Log.Infof("Starting in %s environment", Configuration.GetRunningEnvironment())

	username, _ := ioutil.ReadFile("/run/secrets/basic-auth-user")
	OpenFaaSUsername = strings.TrimSpace(string(username))
	password, _ := ioutil.ReadFile("/run/secrets/basic-auth-password")
	OpenFaaSPassword = strings.TrimSpace(string(password))
	// }

	log.Log.Debug("Init with user %s and password %s", OpenFaaSUsername, OpenFaaSPassword)
	log.Log.Infof("Init with RunningFunctionsMax %d and QueueMaxLength %d", Configuration.GetRunningFunctionMax(), Configuration.GetQueueLengthMax())
}

func Start() {

}
