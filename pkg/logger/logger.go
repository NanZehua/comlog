/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logger

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Level type
type Level uint32

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logger.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

type csiLog struct {
	err        error
	logger     *log.Entry
	target     string
	source     string
	volumeID   string
	snapshotID string
	methodName string
	options    interface{}
}

var _ CSILogger = &csiLog{}

type Fields map[string]interface{}

type CSILogger interface {
	Logf(level Level, format string, args ...interface{})
	WithField(key string, value interface{}) CSILogger
	WithFields(fields Fields) CSILogger
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
}

func New(methodName, nodeID string) CSILogger {

	return &csiLog{
		methodName: methodName,
		logger:     log.New().WithFields(log.Fields{"nodeID": nodeID}),
	}
}

func (csi *csiLog) Logf(level Level, format string, args ...interface{}) {
	// format message information to output
	format = fmt.Sprintf("%s: %s", csi.methodName, format)

	var logLevel log.Level
	switch level {
	case PanicLevel:
		logLevel = log.PanicLevel
	case FatalLevel:
		logLevel = log.FatalLevel
	case ErrorLevel:
		logLevel = log.ErrorLevel
	case WarnLevel:
		logLevel = log.WarnLevel
	case InfoLevel:
		logLevel = log.InfoLevel
	case DebugLevel:
		logLevel = log.WarnLevel
	case TraceLevel:
		logLevel = log.TraceLevel
	default:
		logLevel = log.InfoLevel
	}

	csi.logger.Logf(logLevel, format, args)
}

func (csi *csiLog) WithField(key string, value interface{}) CSILogger {
	csi.logger = csi.logger.WithField(key, value)
	return csi
}

func (csi *csiLog) WithFields(fields Fields) CSILogger {
	var logFileds = make(log.Fields)

	for key, value := range fields {
		logFileds[key] = value
	}

	csi.logger = csi.logger.WithFields(logFileds)

	return csi
}

func (csi *csiLog) Infof(format string, args ...interface{}) {
	csi.Logf(InfoLevel, format, args)
}

func (csi *csiLog) Errorf(format string, args ...interface{}) {
	csi.Logf(ErrorLevel, format, args)
}

func (csi *csiLog) Warningf(format string, args ...interface{}) {
	csi.Logf(WarnLevel, format, args)
}
