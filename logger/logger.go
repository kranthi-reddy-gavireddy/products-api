package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

/*
Common log messages
*/
const (
	invalidArgMessage      = "invalid argument"
	invalidArgValueMessage = "invalid argument value"
	missingArgMessage      = "missing argument"
)

/*
Logger interface
Allows mocking & decoupling from zerolog
*/
type ZeroLogger interface {
	Error() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Debug() *zerolog.Event
	With() zerolog.Context
}

var loggerInstance ZeroLogger

// init runs once when the package is loaded
func init() {
	setDefaultLogLevel()
}

/*
Set log level based on environment
ENV=prod  -> Info
ENV!=prod -> Debug
*/
func setDefaultLogLevel() {
	if os.Getenv("ENV") == "prod" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

/*
Allow injecting a custom logger (useful for tests)
*/
func SetLogger(baseLogger ZeroLogger) {
	loggerInstance = baseLogger
}

/*
Return logger instance (lazy init)
*/
func logger() ZeroLogger {
	if loggerInstance != nil {
		return loggerInstance
	}

	zerolog.TimeFieldFormat = time.RFC1123Z
	baseLogger := baseLogger()
	SetLogger(&baseLogger)
	return loggerInstance
}

func baseLogger() zerolog.Logger {
	return zerolog.New(os.Stderr).
		With().
		Str("service", "products-api").
		Str("env", os.Getenv("ENV")).
		Str("version", os.Getenv("APP_VERSION")).
		Timestamp().
		Logger()
}

/*
	Error helpers
*/

func InvalidArg(argName string) {
	logger().
		Error().
		Str("arg", argName).
		Msg(invalidArgMessage)
}

func WithCorrelationID(cid string) {
	logger := baseLogger().With().Str("cid", cid).Logger()
	SetLogger(&logger)
}

func InvalidArgValue(argName string, argValue interface{}) {
	logger().
		Error().
		Str("arg", argName).
		Interface("value", argValue).
		Msg(invalidArgValueMessage)
}

func MissingArg(argName string) {
	logger().
		Error().
		Str("arg", argName).
		Msg(missingArgMessage)
}

func Error(err error, msg string) {
	logger().
		Error().
		Err(err).
		Msg(msg)
}

func ErrorWithFields(err error, msg string, fields map[string]interface{}) {
	log := logger().With().Fields(fields).Logger()
	log.Error().Err(err).Msg(msg)
}

/*
	Info helpers
*/

func Infof(msg string, v ...interface{}) {
	logger().Info().Msgf(msg, v...)
}

func InfoWithFields(msg string, fields map[string]interface{}) {
	log := logger().With().Fields(fields).Logger()
	log.Info().Msg(msg)
}

/*
	Warn helpers
*/

func Warnf(msg string, v ...interface{}) {
	logger().Warn().Msgf(msg, v...)
}

func WarnWithFields(msg string, fields map[string]interface{}) {
	log := logger().With().Fields(fields).Logger()
	log.Warn().Msg(msg)
}

/*
	Debug helpers
*/

func Debugf(msg string, v ...interface{}) {
	logger().Debug().Msgf(msg, v...)
}

func DebugWithFields(msg string, fields map[string]interface{}) {
	log := logger().With().Fields(fields).Logger()
	log.Debug().Msg(msg)
}

func Errorf(msg string, v ...interface{}) {
	logger().Error().Msgf(msg, v...)
}
