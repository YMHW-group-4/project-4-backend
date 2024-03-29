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
	Seed     string
}

// getConfigFromEnv retrieves configuration from the environment, if environment
// is not set then default configuration will be given instead.
func getConfigFromEnv() Configuration {
	interval := util.GetEnv("INTERVAL", "20m")

	if _, err := time.ParseDuration(interval); err != nil {
		interval = "20m"
	}

	return Configuration{
		Debug:    util.GetEnv("DEBUG", false),
		Port:     util.GetEnv("PORT", 30333),
		APIPort:  util.GetEnv("API_PORT", 8080),
		Interval: interval,
		Seed:     util.GetEnv("DNS_SEED", "localhost:3000"),
	}
}
