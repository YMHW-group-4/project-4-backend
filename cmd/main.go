package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var node *Node

// setupLogger checks whether the Stdout is a terminal or not
// if so it sets the global log's writer to a ConsoleWriter.
func setupLogger() {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		logOutput := zerolog.ConsoleWriter{Out: os.Stderr}
		log.Logger = log.Output(logOutput)

		return
	}

	log.Logger = log.Output(os.Stderr)
}

// setLogLevel sets the global log level to either debug or info level.
func setLogLevel(debug bool) {
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func main() {
	startup := time.Now()
	config := getConfigFromEnv()
	setLogLevel(config.Debug)
	setupLogger()

	log.Info().
		Str("version", version).
		Int("port", config.Port).
		Int("api", config.APIPort).
		Str("interval", config.Interval).
		Bool("debug", config.Debug).
		Msg("node: startup")

	n, err := NewNode(config)
	if err != nil {
		log.Fatal().Err(err).Msg("node: failed to create node")
	}

	node = n

	api := NewAPI(config.APIPort, config.Seed)

	node.Run()

	log.Info().
		TimeDiff("startup", node.Uptime, startup).
		Msg("node: running")

	api.Start()

	node.HandleSigterm()

	api.Stop()
}
