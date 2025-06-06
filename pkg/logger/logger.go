// logger helper 
package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	configureLogger()
}

// configureLogger sets up the global Log instance based on environment variables.
// - LOG_LEVEL: e.g. "debug", "info", "warn", "error" (defaults to "info")
// - LOG_FORMAT: "json" or "text" (defaults to "text")
func configureLogger() {
	// 1) Set log level
	levelStr := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if levelStr == "" {
		levelStr = "info"
	}
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		// Fallback to Info if invalid
		level = logrus.InfoLevel
	}
	Log.SetLevel(level)

	// 2) Set formatter
	format := strings.ToLower(os.Getenv("LOG_FORMAT"))
	switch format {
	case "json":
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		// TextFormatter with full timestamp
		Log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat:  "2006-01-02 15:04:05",
			FullTimestamp:    true,
			DisableColors:    false,
			QuoteEmptyFields: true,
		})
	}

	// 3) Output (default is os.Stdout)
	Log.SetOutput(os.Stdout)
}

// Below are convenient wrapper methods. You can also use Log.WithFields(...) directly.

func Debug(args ...interface{}) {
	Log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

func Warn(args ...interface{}) {
	Log.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	Log.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	Log.Panicf(format, args...)
}
