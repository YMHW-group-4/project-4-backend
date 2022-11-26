package main

import (
	"backend/API"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// setupLogger checks whether the Stdout is a terminal or not.
// if so it sets the global log's writer to a ConsoleWriter.
func setupLogger() {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		logOutput := zerolog.ConsoleWriter{Out: os.Stderr}
		log.Logger = log.Output(logOutput)

		return
	}

	log.Logger = log.Output(os.Stderr)
}

// setLogLevel sets the global logging level of the log's writer.
func setLogLevel(level zerolog.Level) {
	if zerolog.GlobalLevel() != level {
		zerolog.SetGlobalLevel(level)
	}
}

func main() {
	API.StartServer()
	setupLogger()
	setLogLevel(zerolog.DebugLevel)

	log.Info().Msgf("%s", greet("John"))
}

func greet(s string) string {
	return fmt.Sprintf("Hello %s!", s)
}
