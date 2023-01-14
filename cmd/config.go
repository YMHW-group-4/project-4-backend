package main

import "backend/util"

// Configuration all configuration required by the node.
type Configuration struct {
	Debug   bool
	Port    int
	ApiPort int
}

// getConfigFromEnv retrieves configuration from the environment, if environment
// is not set then default configuration will be given instead.
func getConfigFromEnv() Configuration {
	return Configuration{
		Debug:   util.GetEnv("DEBUG", false),
		Port:    util.GetEnv("PORT", 30333),    //nolint
		ApiPort: util.GetEnv("API_PORT", 8080), //nolint
	}
}
