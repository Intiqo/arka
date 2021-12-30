package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Bootstrap() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	if os.Getenv("APP_PRODUCTION") == "true" {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("***%s****", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	Log = zerolog.New(output).With().Timestamp().Caller().Logger()

	Log.Debug().Msg("initialized logger...")
}
