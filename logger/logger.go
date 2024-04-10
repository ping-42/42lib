package logger

import (
	"github.com/ping-42/42lib/config"
	log "github.com/sirupsen/logrus"
)

const (
// testType is used as a JSON key to indicate the testing context
// logKeyTestType = "testType"
)

// Logger is defined during initialization and kept globally to be used everywhere
var Logger *log.Logger

func init() {
	Logger = log.New()
	if config.GetConfig().Env != config.Dev {
		Logger.SetFormatter(&log.JSONFormatter{})
	}
}

// Base will return a *log.entry with set service name JSON key as value
// If predefined log entry is supplied, the JSON fields will be added on top of that
func Base(service string, predefined ...*log.Entry) *log.Entry {
	testTypeField := log.Fields{
		"service": service,
	}
	if len(predefined) > 0 {
		return predefined[0].WithFields(testTypeField)
	}

	return Logger.WithFields(testTypeField)
}

// WithTestType will return a *log.entry with set test type JSON key as value
// If predefined log entry is supplied, the JSON fields will be added on top of that
func WithTestType(tt string, predefined ...*log.Entry) *log.Entry {
	testTypeField := log.Fields{
		"testType": tt,
	}
	if len(predefined) > 0 {
		return predefined[0].WithFields(testTypeField)
	}

	return Logger.WithFields(testTypeField)
}

// LogError logs error with "msg" and sets the "err" as value of error key
func LogError(err, msg string, predefined *log.Entry) {
	if predefined == nil {
		Logger.WithFields(log.Fields{
			"error": err,
			"msg":   msg,
		}).Error("logger.go/LogError: called with nil log Entry")
		return
	}

	errField := log.Fields{
		"error": err,
	}
	predefined.WithFields(errField).Error(msg)
}
