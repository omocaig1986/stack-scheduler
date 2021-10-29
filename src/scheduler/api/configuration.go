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
	"io/ioutil"
	"net/http"
	"scheduler/config"
	"scheduler/errors"
	"scheduler/log"
	"scheduler/scheduler"
	"scheduler/types"
	"scheduler/utils"
)

// Retrieve the current configuration of the system.
func GetConfiguration(w http.ResponseWriter, r *http.Request) {
	configuration, err := json.Marshal(config.Configuration.GetConfiguration())
	if err != nil {
		log.Log.Errorf("Cannot encode configuration to json")
		errors.ReplyWithError(w, errors.GenericError)
		return
	}

	utils.SendJSONResponse(&w, 200, string(configuration))
}

// Set the configuration of the system (conform to config.ConfigurationSetExp) and save it to a file, in such a way it is
// load at the startup. This configuration does not include the scheduler information.
func SetConfiguration(w http.ResponseWriter, r *http.Request) {
	defaultConfiguration := config.GetDefaultConfiguration()
	currentConfiguration := config.Configuration.GetConfiguration()
	reqBody, _ := ioutil.ReadAll(r.Body)

	var newConfiguration *config.ConfigurationSetExp
	var err error
	// do the merge with the default configuration or existing
	if config.ConfigurationReadFromFile {
		err = json.Unmarshal(reqBody, &currentConfiguration)
		newConfiguration = currentConfiguration
	} else {
		err = json.Unmarshal(reqBody, &defaultConfiguration)
		newConfiguration = defaultConfiguration
	}
	if err != nil {
		log.Log.Errorf("Cannot encode passed configuration")
		errors.ReplyWithError(w, errors.GenericError)
		return
	}

	// update existing configuration
	config.Configuration.SetConfiguration(newConfiguration)

	// save configuration to file
	err = config.SaveConfigurationToConfigFile()
	if err != nil {
		log.Log.Warningf("Cannot save configuration to file %s", config.GetConfigFilePath())
	}

	log.Log.Infof("Configuration updated")

	w.WriteHeader(200)
}

// Retrieve the scheduler information.
func GetScheduler(w http.ResponseWriter, r *http.Request) {
	sched, err := json.Marshal(scheduler.GetScheduler())
	if err != nil {
		log.Log.Errorf("Cannot encode configuration to json")
		errors.ReplyWithError(w, errors.GenericError)
		return
	}

	utils.SendJSONResponse(&w, 200, string(sched))
}

// Set the scheduler information and save the configuration to file in such a way it is loaded automatically at startup.
func SetScheduler(w http.ResponseWriter, r *http.Request) {
	var proposedScheduler = types.SchedulerDescriptor{}
	reqBody, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(reqBody, &proposedScheduler)
	if err != nil {
		log.Log.Errorf("Cannot decode passed configuration: %s", err.Error())
		errors.ReplyWithError(w, errors.InputNotValid)
		return
	}

	err = scheduler.SetScheduler(&proposedScheduler)
	if err != nil {
		log.Log.Errorf("Cannot set new scheduler: %s", err.Error())
		errors.ReplyWithErrorMessage(w, errors.GenericError, err.Error())
		return
	}

	// save configuration to file
	err = config.SaveConfigurationSchedulerToConfigFile(scheduler.GetScheduler())
	if err != nil {
		log.Log.Errorf("Cannot save configuration to file %s", config.GetConfigSchedulerFilePath())
	}

	log.Log.Infof("Configuration updated with scheduler: %s", scheduler.GetName())

	w.WriteHeader(200)
}
