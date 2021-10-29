package config

import "os"

// GetDataPath returns the path in which the application can write data files
func GetDataPath() string {
	if os.Getenv(EnvRunningEnvironment) != RunningEnvironmentProduction {
		return "." + DataPath
	}
	return DataPath
}
