package logrus

import (
	"io"
)

// Methods required for a logger.
type Log interface {
	Fatalf(message string, args ...interface{})
	Panicf(message string, args ...interface{})
	Debugf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	LogLevel() string
	SetLevel(string)
	Tracef(message string, args ...interface{})
	Warningf(message string, args ...interface{})

	WithError(error) Log
	WithField(string, interface{}) Log
	WithFields(map[string]interface{}) Log

	SetOutput(io.Writer)
}
