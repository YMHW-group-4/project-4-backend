package main

import (
	"time"

	"backend/util"
)

// Configuration all configuration required by the node.
type Configuration struct {
	Debug    bool
	Port     int
	APIPort  int
	Interval string
}

// getConfigFromEnv retrieves configuration from the environment, if environment
// is not set then default configuration will be given instead.
func getConfigFromEnv() Configuration {
	interval := util.GetEnv("INTERVAL", "10m")

	_, err := time.ParseDuration(interval)
	if err != nil {
		interval = "10m"
	}

	return Configuration{
		Debug:    util.GetEnv("DEBUG", false),
		Port:     util.GetEnv("PORT", 30333),    //nolint
		APIPort:  util.GetEnv("API_PORT", 8080), //nolint
		Interval: interval,
	}
}
