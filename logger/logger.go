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

type BaseLogger struct {
	zerolog.Logger
}

/*
Logger interface
Allows mocking & decoupling from zerolog
*/
type ZeroLogger interface {
	Infof(msg string, v ...interface{})
	Errorf(msg string, v ...interface{})
	Warnf(msg string, v ...interface{})
	Debugf(msg string, v ...interface{})
	InfoWithFields(msg string, fields map[string]interface{})
	ErrorWithFields(err error, msg string, fields map[string]interface{})
	ErrorMsg(err error, msg string)
	WarnWithFields(msg string, fields map[string]interface{})
	DebugWithFields(msg string, fields map[string]interface{})
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
func SetLogger(baseLogger ZeroLogger) ZeroLogger {
	loggerInstance = baseLogger
	return loggerInstance
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
	//SetLogger(&baseLogger)
	return &BaseLogger{Logger: baseLogger}
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

func WithCorrelationID(cid string) ZeroLogger {
	logger := baseLogger().With().Str("cid", cid).Logger()
	return &BaseLogger{Logger: logger}
}

func InvalidArgValue(argName string, argValue interface{}) {
	logger().
		Error().
		Str("arg", argName).
		Interface("value", argValue).
		Msg(invalidArgValueMessage)
}

func (b *BaseLogger) MissingArg(argName string) {
	b.Logger.Error().Str("arg", argName).Msg(missingArgMessage)
}

func (b *BaseLogger) ErrorMsg(err error, msg string) {
	b.Logger.Error().Err(err).Msg(msg)
}

func (b *BaseLogger) ErrorWithFields(err error, msg string, fields map[string]interface{}) {
	log := b.With().Fields(fields).Logger()
	log.Error().Err(err).Msg(msg)
}

/*
	Info helpers
*/

func (b *BaseLogger) Infof(msg string, v ...interface{}) {
	b.Logger.Info().Msgf(msg, v...)
}

func (b *BaseLogger) InfoWithFields(msg string, fields map[string]interface{}) {
	log := b.With().Fields(fields).Logger()
	log.Info().Msg(msg)
}

/*
	Warn helpers
*/

func (b *BaseLogger) Warnf(msg string, v ...interface{}) {
	b.Logger.Warn().Msgf(msg, v...)
}

func (b *BaseLogger) WarnWithFields(msg string, fields map[string]interface{}) {
	log := b.With().Fields(fields).Logger()
	log.Warn().Msg(msg)
}

/*
	Debug helpers
*/

func (b *BaseLogger) Debugf(msg string, v ...interface{}) {
	b.Logger.Debug().Msgf(msg, v...)
}

func (b *BaseLogger) DebugWithFields(msg string, fields map[string]interface{}) {
	log := b.With().Fields(fields).Logger()
	log.Debug().Msg(msg)
}

func (b *BaseLogger) Errorf(msg string, v ...interface{}) {
	b.Logger.Error().Msgf(msg, v...)
}

func New() ZeroLogger {
	logger := baseLogger()
	return &BaseLogger{Logger: logger}
}
