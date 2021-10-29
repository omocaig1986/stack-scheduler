package config

import (
	"encoding/json"
	"io/ioutil"
	"scheduler/log"
	"scheduler/types"
)

func GetConfigFilePath() string {
	return GetDataPath() + "/" + ConfigurationFileName
}

func GetConfigSchedulerFilePath() string {
	return GetDataPath() + "/" + ConfigurationSchedulerFileName
}

func SaveConfigurationToConfigFile() error {
	// prepare configuration
	confExported := GetDefaultConfiguration()
	copyAllFieldsToExp(&Configuration, confExported)

	// save configuration to file
	configJson, err := json.MarshalIndent(confExported, "", "  ")
	err = ioutil.WriteFile(GetConfigFilePath(), configJson, 0644)
	if err != nil {
		log.Log.Errorf("Cannot save configuration to file %s: %s", GetConfigFilePath(), err.Error())
		return err
	}

	return nil
}

func SaveConfigurationSchedulerToConfigFile(descriptor *types.SchedulerDescriptor) error {
	configJson, err := json.MarshalIndent(descriptor, "", "  ")
	err = ioutil.WriteFile(GetConfigSchedulerFilePath(), configJson, 0644)
	if err != nil {
		log.Log.Errorf("Cannot save configuration to file %s: %s", GetConfigSchedulerFilePath(), err.Error())
		return err
	}

	return nil
}
