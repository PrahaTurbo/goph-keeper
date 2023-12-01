package logger

import (
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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
