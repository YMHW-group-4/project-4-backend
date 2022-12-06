package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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
	setLogLevel(config.debug)
	setupLogger()

	log.Info().
		Str("version", version).
		Int("port", config.port).
		Bool("debug", config.debug).
		Msg("node: startup")

	node, err := newNode(config)
	if err != nil {
		log.Fatal().Err(err).Msg("node: failed to create node")
	}

	node.run()

	log.Info().
		TimeDiff("startup", time.Now(), startup).
		Msg("node: running")

	node.handleSigterm()
}
