// Package logger provides methods to create and customize the application's logger.
package logger

import (
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// NewLogger constructs a new zerolog.Logger with a custom caller marshal function.
// The logger includes the caller (file and line number) of log events.
func NewLogger() zerolog.Logger {
	zerolog.CallerMarshalFunc = customCallerMarshal

	logger := log.With().Caller().Logger()

	return logger
}

func customCallerMarshal(pc uintptr, file string, line int) string {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err)
	}

	filePath := strings.ReplaceAll(file, root, "")

	return filePath + ":" + strconv.Itoa(line)
}
