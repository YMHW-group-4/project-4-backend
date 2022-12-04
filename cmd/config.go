package main

import "backend/util"

// configuration all configuration required by the node.
type configuration struct {
	debug bool
	port  int
}

// getConfigFromEnv retrieves configuration from the environment, if environment
// is not set then default configuration will be given instead.
func getConfigFromEnv() configuration {
	return configuration{
		debug: util.GetEnv("NODE_DEBUG", false),
		port:  util.GetEnv("NODE_PORT", 30333), //nolint
	}
}
