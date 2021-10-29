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
	"encoding/json"
	"io/ioutil"
	"os"
	"scheduler/log"
)

type ConfigError struct{}

func (ConfigError) Error() string {
	return "Configuration Error"
}

type ConfigurationSet struct {
	runningFunctionMax     uint
	queueLengthMax         uint
	listeningPort          uint
	openFaasListeningPort  uint
	openFaasListeningHost  string
	discoveryListeningPort uint
	discoveryListeningHost string
	runningEnvironment     string
}

type ConfigurationSetExp struct {
	RunningFunctionMax     uint   `json:"running_functions_max" bson:"running_functions_max"`
	QueueLengthMax         uint   `json:"queue_length_max" bson:"queue_length_max"`
	ListeningPort          uint   `json:"listening_port" bson:"listening_port"`
	OpenFaasListeningPort  uint   `json:"faas_listening_port" bson:"faas_listening_port"`
	OpenFaasListeningHost  string `json:"faas_listening_host" bson:"faas_listening_host"`
	DiscoveryListeningPort uint   `json:"discovery_listening_port" bson:"discovery_listening_port"`
	DiscoveryListeningHost string `json:"discovery_listening_host" bson:"discovery_listening_host"`
	RunningEnvironment     string `json:"running_environment" bson:"running_environment"`
}

/*
 * Getters
 */

func (c ConfigurationSet) GetRunningFunctionMax() uint {
	return c.runningFunctionMax
}
func (c ConfigurationSet) GetQueueLengthMax() uint {
	return c.queueLengthMax
}
func (c ConfigurationSet) GetListeningPort() uint {
	return c.listeningPort
}
func (c ConfigurationSet) GetOpenFaasListeningPort() uint {
	return c.openFaasListeningPort
}
func (c ConfigurationSet) GetOpenFaasListeningHost() string {
	return c.openFaasListeningHost
}
func (c ConfigurationSet) GetDiscoveryListeningPort() uint {
	return c.discoveryListeningPort
}
func (c ConfigurationSet) GetDiscoveryListeningHost() string {
	if os.Getenv(EnvRunningEnvironment) != RunningEnvironmentProduction {
		return "localhost"
	}
	return c.discoveryListeningHost
}
func (c ConfigurationSet) GetRunningEnvironment() string {
	return c.runningEnvironment
}

// GetConfiguration obtains the configuration with exported fields
func (c ConfigurationSet) GetConfiguration() *ConfigurationSetExp {
	conf := &ConfigurationSetExp{}
	copyAllFieldsToExp(&Configuration, conf)
	return conf
}

/*
 * Setters
 */

func (c *ConfigurationSet) SetRunningFunctionMax(n uint) {
	c.runningFunctionMax = n
}
func (c *ConfigurationSet) SetQueueLengthMax(n uint) {
	c.queueLengthMax = n
}
func (c *ConfigurationSet) SetListeningPort(n uint) {
	c.listeningPort = n
}
func (c *ConfigurationSet) SetFaasListeningPort(n uint) {
	c.openFaasListeningPort = n
}
func (c *ConfigurationSet) SetFaasListeningHost(s string) {
	c.openFaasListeningHost = s
}
func (c *ConfigurationSet) SetDiscoveryListeningPort(n uint) {
	c.discoveryListeningPort = n
}
func (c *ConfigurationSet) SetDiscoveryListeningHost(s string) {
	c.discoveryListeningHost = s
}
func (c *ConfigurationSet) SetRunningEnvironment(s string) {
	c.runningEnvironment = s
}

// SetConfiguration updates the entire configuration
func (c *ConfigurationSet) SetConfiguration(exp *ConfigurationSetExp) {
	copyAllFieldsToUnExp(exp, c)
}

/*
 * Utils
 */

func ReadConfigFile() (*ConfigurationSet, error) {
	file, err := ioutil.ReadFile(GetConfigFilePath())

	conf := GetDefaultConfiguration()
	confValid := ConfigurationSet{}

	if err != nil {
		log.Log.Info("Cannot read configuration file at %s", GetConfigFilePath())
	} else {
		err = json.Unmarshal(file, &conf)
		if err != nil {
			log.Log.Errorf("Cannot decode configuration file, maybe not valid json: %s", err.Error())
		}
	}

	// get running env always from env, if not specified we assume production
	if os.Getenv(EnvRunningEnvironment) == RunningEnvironmentDevelopment {
		conf.RunningEnvironment = RunningEnvironmentDevelopment
	} else {
		conf.RunningEnvironment = RunningEnvironmentProduction
	}

	copyAllFieldsToUnExp(conf, &confValid)

	// update config field
	Configuration = confValid

	return &confValid, nil
}

func GetDefaultConfiguration() *ConfigurationSetExp {
	return &ConfigurationSetExp{
		RunningFunctionMax:     DefaultFunctionsRunningMax,
		QueueLengthMax:         DefaultQueueLengthMax,
		ListeningPort:          DefaultListeningPort,
		OpenFaasListeningPort:  DefaultOpenFaaSListeningPort,
		OpenFaasListeningHost:  DefaultOpenFaaSListeningHost,
		DiscoveryListeningPort: DefaultStackDiscoveryListeningPort,
		DiscoveryListeningHost: DefaultStackDiscoveryListeningHost,
		RunningEnvironment:     os.Getenv(EnvRunningEnvironment),
	}
}

func copyAllFieldsToExp(from *ConfigurationSet, to *ConfigurationSetExp) {
	to.RunningFunctionMax = from.runningFunctionMax
	to.QueueLengthMax = from.queueLengthMax
	to.ListeningPort = from.listeningPort
	to.OpenFaasListeningHost = from.openFaasListeningHost
	to.OpenFaasListeningPort = from.openFaasListeningPort
	to.DiscoveryListeningHost = from.discoveryListeningHost
	to.DiscoveryListeningPort = from.discoveryListeningPort
	to.RunningEnvironment = from.runningEnvironment
}

func copyAllFieldsToUnExp(from *ConfigurationSetExp, to *ConfigurationSet) {
	to.runningFunctionMax = from.RunningFunctionMax
	to.queueLengthMax = from.QueueLengthMax
	to.listeningPort = from.ListeningPort
	to.openFaasListeningHost = from.OpenFaasListeningHost
	to.openFaasListeningPort = from.OpenFaasListeningPort
	to.discoveryListeningHost = from.DiscoveryListeningHost
	to.discoveryListeningPort = from.DiscoveryListeningPort
	to.runningEnvironment = from.RunningEnvironment
}
